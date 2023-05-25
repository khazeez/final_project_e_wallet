package handler

import (
	
	"net/http"
	"strconv"
	
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
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

func (h *PaymentHandler) MakePayment(c *gin.Context) {
	var payment domain.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.paymentUsecase.MakePayment(&payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment made successfully"})
}
