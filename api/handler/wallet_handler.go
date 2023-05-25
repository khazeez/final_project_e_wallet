package handler

import (
	
	"net/http"
	"strconv"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)

type WalletHandler struct {
	walletUsecase app.WalletUsecase
}

func NewWalletHandler(walletUsecase app.WalletUsecase) *WalletHandler {
	return &WalletHandler{
		walletUsecase: walletUsecase,
	}
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var wallet domain.Wallet
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.walletUsecase.CreateWallet(&wallet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Wallet created successfully"})
}

func (h *WalletHandler) GetWalletByID(c *gin.Context) {
	walletID, err := strconv.Atoi(c.Param("walletID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	wallet, err := h.walletUsecase.GetWalletByID(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) UpdateWalletBalance(c *gin.Context) {
	walletID, err := strconv.Atoi(c.Param("walletID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	var wallet domain.Wallet
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wallet.ID = walletID

	if err := h.walletUsecase.UpdateWalletBalanceUpdate(&wallet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet balance updated successfully"})
}

func (h *WalletHandler) DeleteWallet(c *gin.Context) {
	walletID, err := strconv.Atoi(c.Param("walletID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	if err := h.walletUsecase.DeleteWallet(walletID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted successfully"})
}
