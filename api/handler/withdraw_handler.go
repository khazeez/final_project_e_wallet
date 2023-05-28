package handler

import (
	"net/http"
	"strconv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/jung-kurt/gofpdf"
	"bytes"

)

type WithdrawalHandler struct {
	withdrawalUsecase app.WithdrawUsecase
}

func NewWithdrawalHandler(withdrawalUsecase app.WithdrawUsecase) *WithdrawalHandler {
	return &WithdrawalHandler{
		withdrawalUsecase: withdrawalUsecase,
	}
}

func (h *WithdrawalHandler) CreateWithdrawal(c *gin.Context) {
	var withdrawal domain.Withdrawal
	if err := c.ShouldBindJSON(&withdrawal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.withdrawalUsecase.CreateWithdrawal(&withdrawal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Withdrawal created successfully"})
}

func (h *WithdrawalHandler) GetWithdrawalByID(c *gin.Context) {
	withdrawalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid withdrawal ID"})
		return
	}

	withdrawal, err := h.withdrawalUsecase.GetWithdrawalByID(withdrawalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, withdrawal)
}

func (h *WithdrawalHandler) UpdateWithdrawal(c *gin.Context) {
	withdrawalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid withdrawal ID"})
		return
	}

	var withdrawal domain.Withdrawal
	if err := c.ShouldBindJSON(&withdrawal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	withdrawal.ID = withdrawalID

	if err := h.withdrawalUsecase.UpdateWithdrawal(&withdrawal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal updated successfully"})
}

func (h *WithdrawalHandler) DeleteWithdrawal(c *gin.Context) {
	withdrawalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid withdrawal ID"})
		return
	}

	if err := h.withdrawalUsecase.DeleteWithdrawal(withdrawalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal deleted successfully"})
}

func (h *WithdrawalHandler) MakeWithdrawal(c *gin.Context) {
	var withdrawal domain.Withdrawal
	if err := c.ShouldBindJSON(&withdrawal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.withdrawalUsecase.MakeWithdrawal(&withdrawal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal made successfully"})
}
func (h *WithdrawalHandler) HistoryTransaction(c *gin.Context) {
	walletID, err := strconv.Atoi(c.Param("wallet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	withdrawals, err := h.withdrawalUsecase.HistoryTransaction(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate PDF from transaction data
	pdfOutput := GeneratePDFwd(withdrawals)

	// Send PDF file as response
	c.Data(http.StatusOK, "application/pdf", pdfOutput)
}

func GeneratePDFwd(withdrawals []*domain.Withdrawal) []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Transaction History")

	// Add transaction data to PDF
	for _, withdrawal := range withdrawals {
		pdf.Ln(12)
		pdf.Cell(40, 10, fmt.Sprintf("Withdrawal ID: %d", withdrawal.ID))
		pdf.Cell(40, 10, fmt.Sprintf("Amount: %f", withdrawal.Amount))
	}
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
