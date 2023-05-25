package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)

type WithdrawHandler struct {
	withdrawUsecase app.WithdrawUsecase
}

func NewWithdrawHandler(withdrawUsecase app.WithdrawUsecase) *WithdrawHandler {
	return &WithdrawHandler{
		withdrawUsecase: withdrawUsecase,
	}
}

func (h *WithdrawHandler) CreateWithdrawal(w http.ResponseWriter, r *http.Request) {
	var withdrawal domain.Withdrawal
	err := json.NewDecoder(r.Body).Decode(&withdrawal)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.withdrawUsecase.MakeWithdrawal(&withdrawal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *WithdrawHandler) GetWithdrawalByID(w http.ResponseWriter, r *http.Request) {
	withdrawalIDStr := r.URL.Query().Get("withdrawal_id")
	withdrawalID, err := strconv.Atoi(withdrawalIDStr)
	if err != nil {
		http.Error(w, "Invalid withdrawal ID", http.StatusBadRequest)
		return
	}

	withdrawal, err := h.withdrawUsecase.GetWithdrawalByID(withdrawalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(withdrawal)
}

func (h *WithdrawHandler) UpdateWithdrawal(w http.ResponseWriter, r *http.Request) {
	var withdrawal domain.Withdrawal
	err := json.NewDecoder(r.Body).Decode(&withdrawal)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.withdrawUsecase.UpdateWithdrawal(&withdrawal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WithdrawHandler) DeleteWithdrawal(w http.ResponseWriter, r *http.Request) {
	withdrawalIDStr := r.URL.Query().Get("withdrawal_id")
	withdrawalID, err := strconv.Atoi(withdrawalIDStr)
	if err != nil {
		http.Error(w, "Invalid withdrawal ID", http.StatusBadRequest)
		return
	}

	err = h.withdrawUsecase.DeleteWithdrawal(withdrawalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
