package app

import (
	// "fmt"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
)

type TransferUsecase interface {
	CreateTransfer(transfer *domain.Transfer) error
	GetTransferByID(transferID int) (*domain.Transfer, error)
	UpdateTransfer(transfer *domain.Transfer) error
	DeleteTransfer(transferID int) error
	HistoryTransaction(transferID int) ([]*domain.Transfer, error)
	// MakeTransfer(transfer *domain.Transfer) error
}

type transferUsecase struct {
	transferRepository repository.TransferRepository
	walletRepository  repository.WalletRepository
}

func NewTransferUsecase(transferRepository repository.TransferRepository, walletRepository repository.WalletRepository) TransferUsecase {
	return &transferUsecase{
		transferRepository: transferRepository,
		walletRepository:   walletRepository,
	}
}

func (u *transferUsecase) CreateTransfer(transfer *domain.Transfer) error {
	return u.transferRepository.Create(transfer)
}

func (u *transferUsecase) GetTransferByID(transferID int) (*domain.Transfer, error) {
	return u.transferRepository.FindOne(transferID)
}

func (u *transferUsecase) UpdateTransfer(transfer *domain.Transfer) error {
	return u.transferRepository.Update(transfer)
}

func (u *transferUsecase) DeleteTransfer(transferID int) error {
	return u.transferRepository.Delete(transferID)
}
func (u *transferUsecase) HistoryTransaction(transferID int) ([]*domain.Transfer, error) {
	return u.transferRepository.History(transferID)
}

