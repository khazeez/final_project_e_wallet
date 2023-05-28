package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jung-kurt/gofpdf"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUsecase app.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase app.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase: paymentUsecase,
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var payment domain.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.paymentUsecase.CreatePayment(&payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Payment created successfully"})
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	paymentID, err := strconv.Atoi(c.Param("paymentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := h.paymentUsecase.GetPaymentByID(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
	paymentID, err := strconv.Atoi(c.Param("paymentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var payment domain.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payment.ID = paymentID

	if err := h.paymentUsecase.UpdatePayment(&payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment updated successfully"})
}

func (h *PaymentHandler) DeletePayment(c *gin.Context) {
	paymentID, err := strconv.Atoi(c.Param("paymentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	if err := h.paymentUsecase.DeletePayment(paymentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}


func (h *PaymentHandler) HistoryPayment(c *gin.Context) {
	walletID, err := strconv.Atoi(c.Param("paymentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	payments, err := h.paymentUsecase.HistoryTransaction(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate PDF from transaction data
	pdfOutput := GeneratePDFPayment(payments)

	// Send PDF file as response
	c.Data(http.StatusOK, "application/pdf", pdfOutput)
}


func GeneratePDFPayment(payments []*domain.Payment) []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Transaction History")

	// Add transaction data to PDF
	for _, payment := range payments {
		pdf.Ln(12)
		pdf.Cell(20, 10, fmt.Sprintf("Username: %s \n | Email : %s", payment.WalletId.UserId.Name, payment.WalletId.UserId.Email))
		break

	}

	for _, payment := range payments {
		pdf.Ln(12)
		pdf.Cell(20, 10, fmt.Sprintf("Withdrawal ID: %d \n \n | Amount: %f \n | Time: %v ", payment.ID, payment.Amount, payment.Timestamp))

	}
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

