package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type TransferHandler struct {
	transferUsecase app.TransferUsecase
}

func NewTransferHandler(transferUsecase app.TransferUsecase) *TransferHandler {
	return &TransferHandler{
		transferUsecase: transferUsecase,
	}
}

func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	var transfer domain.Transfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.transferUsecase.CreateTransfer(&transfer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transfer created successfully"})
}

func (h *TransferHandler) GetTransferByID(c *gin.Context) {
	transferID, err := strconv.Atoi(c.Param("transferID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	transfer, err := h.transferUsecase.GetTransferByID(transferID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transfer)
}

func (h *TransferHandler) UpdateTransfer(c *gin.Context) {
	transferID, err := strconv.Atoi(c.Param("transferID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	var transfer domain.Transfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transfer.ID = transferID

	if err := h.transferUsecase.UpdateTransfer(&transfer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer updated successfully"})
}

func (h *TransferHandler) DeleteTransfer(c *gin.Context) {
	transferID, err := strconv.Atoi(c.Param("transferID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	if err := h.transferUsecase.DeleteTransfer(transferID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer deleted successfully"})
}

func (h *TransferHandler) HistoryTransaction(c *gin.Context) {
	senderID, err := strconv.Atoi(c.Param("sender_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	transfers, err := h.transferUsecase.HistoryTransaction(senderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate PDF from transaction data
	pdfOutput := GeneratePDFTransfer(transfers)

	// Send PDF file as response
	c.Data(http.StatusOK, "application/pdf", pdfOutput)
}

func GeneratePDFTransfer(transfers []*domain.Transfer) []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Transaction History")

	// Add transaction data to PDF

	for _, transfer := range transfers {
		pdf.Ln(12)
		pdf.Cell(20, 10, fmt.Sprintf("Sender Name: %s \n ", transfer.SenderId.UserId.Sender_Name))
		break

	}

	for _, transfer := range transfers {
		pdf.Ln(12)
		pdf.Cell(20, 10, fmt.Sprintf("To Receiver  ID: %d \n \n | Amount: %f \n | Time: %v ", transfer.ReceiferId.ID, transfer.Amount, transfer.Timestamp))
	}
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
