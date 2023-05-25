package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
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

func (h *TopupHandler) CreateTopup(w http.ResponseWriter, r *http.Request) {
	var topup domain.TopUp
	err := json.NewDecoder(r.Body).Decode(&topup)
	if err != nil {
		h.handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = h.topupUsecase.CreateTopup(&topup)
	if err != nil {
		h.handleError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create top-up: %v", err))
		return
	}

	h.handleSuccess(w, http.StatusCreated, "Top-up created successfully")
}

func (h *TopupHandler) GetTopupByID(w http.ResponseWriter, r *http.Request) {
	// Parse topupID from request parameters
	topupID := parseTopupID(r)
	if topupID == 0 {
		h.handleError(w, http.StatusBadRequest, "Invalid top-up ID")
		return
	}

	topup, err := h.topupUsecase.GetTopupByID(topupID)
	if err != nil {
		h.handleError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get top-up: %v", err))
		return
	}

	h.handleSuccess(w, http.StatusOK, topup)
}

func (h *TopupHandler) UpdateTopup(w http.ResponseWriter, r *http.Request) {
	// Parse topupID from request parameters
	topupID := parseTopupID(r)
	if topupID == 0 {
		h.handleError(w, http.StatusBadRequest, "Invalid top-up ID")
		return
	}

	var topup domain.TopUp
	err := json.NewDecoder(r.Body).Decode(&topup)
	if err != nil {
		h.handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	topup.ID = topupID

	err = h.topupUsecase.UpdateTopup(&topup)
	if err != nil {
		h.handleError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update top-up: %v", err))
		return
	}

	h.handleSuccess(w, http.StatusOK, "Top-up updated successfully")
}

func (h *TopupHandler) DeleteTopup(w http.ResponseWriter, r *http.Request) {
	// Parse topupID from request parameters
	topupID := parseTopupID(r)
	if topupID == 0 {
		h.handleError(w, http.StatusBadRequest, "Invalid top-up ID")
		return
	}

	err := h.topupUsecase.DeleteTopup(topupID)
	if err != nil {
		h.handleError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete top-up: %v", err))
		return
	}

	h.handleSuccess(w, http.StatusOK, "Top-up deleted successfully")
}

func (h *TopupHandler) GetLastTopupAmount(w http.ResponseWriter, r *http.Request) {
	// Parse walletID from request parameters
	walletID := parseWalletID(r)
	if walletID == 0 {
		h.handleError(w, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	amount, err := h.topupUsecase.GetLastTopupAmount(walletID)
	if err != nil {
		h.handleError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get last top-up amount: %v", err))
		return
	}

	response := struct {
		Amount float64 `json:"amount"`
	}{
		Amount: amount,
	}

	h.handleSuccess(w, http.StatusOK, response)
}

func (h *TopupHandler) handleError(w http.ResponseWriter, statusCode int, message string) {
	response := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h *TopupHandler) handleSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func parseTopupID(r *http.Request) int {
	// Parse topupID from request parameters or body, depending on your implementation
	// Return 0 if topupID is not found or invalid
	return 0
}

func parseWalletID(r *http.Request) int {
	// Parse walletID from request parameters or body, depending on your implementation
	// Return 0 if walletID is not found or invalid
	return 0
}
