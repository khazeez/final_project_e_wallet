package app
import (
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
	"errors"
	"fmt"
)
type WalletUsecase interface {
	CreateWallet(wallet *domain.Wallet) error
	GetWalletByID(walletID int) (*domain.Wallet, error)
	UpdateWalletBalanceUpdate(wallet *domain.Wallet) error
	DeleteWallet(walletID int) error
}
type walletUsecase struct {
	walletRepository repository.WalletRepository
	topupRepository  repository.TopupRepository
}
func NewWalletUsecase(walletRepository repository.WalletRepository, topupRepository repository.TopupRepository) WalletUsecase {
	return &walletUsecase{
		walletRepository: walletRepository,
		topupRepository:  topupRepository,
	}
}
func (u *walletUsecase) CreateWallet(wallet *domain.Wallet) error {
	return u.walletRepository.Create(wallet)
}
func (u *walletUsecase) GetWalletByID(walletID int) (*domain.Wallet, error) {
	return u.walletRepository.FindOne(walletID)
}
func (u *walletUsecase) UpdateWalletBalanceUpdate(wallet *domain.Wallet) error {
	wallet, err := u.walletRepository.FindOne(wallet.ID)
	if err != nil {
		return err
	}
	topupAmount, err := u.topupRepository.GetLastTopupAmount(wallet.ID)
	if err != nil {
		return fmt.Errorf("failed to get top-up amount: %v", err)
	}

	if topupAmount == 0 {
		return errors.New("balance amount must be greater than 0")
	}

	wallet.Balance += topupAmount

	err = u.walletRepository.Update(wallet)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %v", err)
	}

	return nil
}
func (u *walletUsecase) DeleteWallet(walletID int) error {
	return u.walletRepository.Delete(walletID)
}
