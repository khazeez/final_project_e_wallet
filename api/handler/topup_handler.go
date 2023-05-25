package handler

import (

	"net/http"
	"strconv"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"

	"github.com/gin-gonic/gin"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
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
