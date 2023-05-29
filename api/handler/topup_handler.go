package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/jung-kurt/gofpdf"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/gin-gonic/gin"
)

type TopupHandler struct {
	topupUsecase app.TopupUsecase
}

func NewTopupHandler(topupUsecase app.TopupUsecase) *TopupHandler {
	return &TopupHandler{
		topupUsecase: topupUsecase,
	}
}

func (h *TopupHandler) CreateTopup(c *gin.Context) {
	var topup domain.TopUp
	if err := c.ShouldBindJSON(&topup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.topupUsecase.CreateTopup(&topup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Topup created successfully"})
}

func (h *TopupHandler) GetTopupByID(c *gin.Context) {
	topupID, err := strconv.Atoi(c.Param("topupID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topup ID"})
		return
	}

	topup, err := h.topupUsecase.GetTopupByID(topupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, topup)
}

func (h *TopupHandler) UpdateTopup(c *gin.Context) {
	topupID, err := strconv.Atoi(c.Param("topupID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topup ID"})
		return
	}

	var topup domain.TopUp
	if err := c.ShouldBindJSON(&topup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	topup.ID = topupID

	if err := h.topupUsecase.UpdateTopup(&topup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Topup updated successfully"})
}

func (h *TopupHandler) DeleteTopup(c *gin.Context) {
	topupID, err := strconv.Atoi(c.Param("topupID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topup ID"})
		return
	}

	if err := h.topupUsecase.DeleteTopup(topupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Topup deleted successfully"})
}

func (h *TopupHandler) GetLastTopupAmount(c *gin.Context) {
	walletID, err := strconv.Atoi(c.Param("walletID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	amount, err := h.topupUsecase.GetLastTopupAmount(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"amount": amount})
}

func (h *TopupHandler) HistoryTopup(c *gin.Context) {
	walletID, err := strconv.Atoi(c.Param("topupID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	topups, err := h.topupUsecase.HistoryTransaction(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate PDF from transaction data
	pdfOutput := GeneratePDFTopup(topups)

	// Send PDF file as response
	c.Data(http.StatusOK, "application/pdf", pdfOutput)
}

func GeneratePDFTopup(topups []*domain.TopUp) []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Transaction History")

	// Add transaction data to PDF
	for _, topup := range topups {
		pdf.Ln(12)
		pdf.Cell(20, 10, fmt.Sprintf("Username: %s \n | Email : %s", topup.WalletId.UserId.Name, topup.WalletId.UserId.Email))
		break

	}

	for _, topup := range topups {
		pdf.Ln(12)
		pdf.Cell(20, 10, fmt.Sprintf("Top-Up ID: %d \n \n | Amount: %2.f \n | Time: %v ", topup.ID, topup.Amount, topup.Timestamp))

	}
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
