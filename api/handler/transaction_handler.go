package handler

import (
	"net/http"
	"strconv"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUsecase app.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase app.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	// Ambil data dari request
	var request struct {
		WalletID int     `json:"wallet_id"`
		Type     string  `json:"type"`
		Amount   float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Membuat transaksi menggunakan use case
	err := h.transactionUsecase.CreateTransaction(request.WalletID, request.Type, request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully"})
}

func (h *TransactionHandler) GetTransactionsByWalletID(c *gin.Context) {
	// Ambil ID wallet dari path parameter
	walletID, err := strconv.Atoi(c.Param("wallet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	// Mendapatkan riwayat transaksi menggunakan use case
	transactions, err := h.transactionUsecase.GetTransactionsByWalletID(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
