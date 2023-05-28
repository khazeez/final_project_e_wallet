package app

import (
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
	"time"
)

type TransactionUsecase interface {
	CreateTransaction(walletID int, transactionType string, amount float64) error
	GetTransactionsByWalletID(walletID int) ([]*domain.Transaction, error)
}

type transactionUsecase struct {
	transactionRepository repository.TransactionRepository
}

func NewTransactionUsecase(transactionRepository repository.TransactionRepository) TransactionUsecase {
	return &transactionUsecase{
		transactionRepository: transactionRepository,
	}
}

func (u *transactionUsecase) CreateTransaction(walletID int, transactionType string, amount float64) error {
	transaction := &domain.Transaction{
		WalletId:  walletID,
		Type:      transactionType,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	err := u.transactionRepository.CreateTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
}

func (u *transactionUsecase) GetTransactionsByWalletID(walletID int) ([]*domain.Transaction, error) {
	transactions, err := u.transactionRepository.GetTransactionsByWalletID(walletID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
