package handler

import (
	"net/http"
	"strconv"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
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

// func (h *TransferHandler) MakeTransfer(c *gin.Context) {
// 	var transfer domain.Transfer
// 	if err := c.ShouldBindJSON(&transfer); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := h.transferUsecase.MakeTransfer(&transfer); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Transfer made successfully"})
// }
