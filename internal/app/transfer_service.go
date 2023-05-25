package app

import (
	"fmt"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
)

type TransferUsecase interface {
	CreateTransfer(transfer *domain.Transfer) error
	GetTransferByID(transferID int) (*domain.Transfer, error)
	UpdateTransfer(transfer *domain.Transfer) error
	DeleteTransfer(transferID int) error
	MakeTransfer(transfer *domain.Transfer) error
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

func (u *transferUsecase) MakeTransfer(transfer *domain.Transfer) error {
	senderWallet, err := u.walletRepository.FindOne(transfer.SenderId.ID)
	if err != nil {
		return fmt.Errorf("failed to find sender wallet: %v", err)
	}

	receiverWallet, err := u.walletRepository.FindOne(transfer.ReceiferId.ID)
	if err != nil {
		return fmt.Errorf("failed to find receiver wallet: %v", err)
	}

	if senderWallet.Balance < transfer.Amount {
		return fmt.Errorf("insufficient balance for transfer")
	}

	senderWallet.Balance -= transfer.Amount
	receiverWallet.Balance += transfer.Amount

	err = u.walletRepository.Update(senderWallet)
	if err != nil {
		return fmt.Errorf("failed to update sender wallet: %v", err)
	}

	err = u.walletRepository.Update(receiverWallet)
	if err != nil {
		return fmt.Errorf("failed to update receiver wallet: %v", err)
	}

	err = u.transferRepository.Create(transfer)
	if err != nil {
		// Rollback the balance changes if failed to create transfer
		senderWallet.Balance += transfer.Amount
		receiverWallet.Balance -= transfer.Amount

		_ = u.walletRepository.Update(senderWallet)
		_ = u.walletRepository.Update(receiverWallet)

		return fmt.Errorf("failed to create transfer: %v", err)
	}

	return nil
}
