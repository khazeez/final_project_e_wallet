package app
import (
	"fmt"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
)
type WithdrawUsecase interface {
	CreateWithdrawal(withdrawal *domain.Withdrawal) error
	GetWithdrawalByID(withdrawalID int) (*domain.Withdrawal, error)
	HistoryTransaction(withdrawalID int) ([]*domain.Withdrawal, error)
	UpdateWithdrawal(withdrawal *domain.Withdrawal) error
	DeleteWithdrawal(withdrawalID int) error
	MakeWithdrawal(withdrawal *domain.Withdrawal) error
}
type withdrawUsecase struct {
	withdrawRepository repository.WithdrawRepository
	walletRepository   repository.WalletRepository
}

func NewWithdrawUsecase(withdrawRepository repository.WithdrawRepository, walletRepository repository.WalletRepository) WithdrawUsecase {
	return &withdrawUsecase{
		withdrawRepository: withdrawRepository,
		walletRepository:   walletRepository,
	}
}

func (u *withdrawUsecase) CreateWithdrawal(withdrawal *domain.Withdrawal) error {
	return u.withdrawRepository.Create(withdrawal)
}

func (u *withdrawUsecase) GetWithdrawalByID(withdrawalID int) (*domain.Withdrawal, error) {
	return u.withdrawRepository.FindOne(withdrawalID)
}

func (u *withdrawUsecase) UpdateWithdrawal(withdrawal *domain.Withdrawal) error {
	return u.withdrawRepository.Update(withdrawal)
}

func (u *withdrawUsecase) DeleteWithdrawal(withdrawalID int) error {
	return u.withdrawRepository.Delete(withdrawalID)
}

func (u *withdrawUsecase) MakeWithdrawal(withdrawal *domain.Withdrawal) error {
	wallet, err := u.walletRepository.FindOne(withdrawal.WalletId.ID)
	if err != nil {
		return fmt.Errorf("failed to find wallet: %v", err)
	}
	if wallet.Balance < float64(withdrawal.Amount) {
		return fmt.Errorf("insufficient balance in wallet")
	}
	// Kurangi saldo pada wallet
	wallet.Balance -= float64(withdrawal.Amount)
	err = u.walletRepository.Update(wallet)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}
	// Simpan withdrawal
	err = u.withdrawRepository.Create(withdrawal)
	if err != nil {
		// Jika gagal menyimpan withdrawal, tambahkan kembali saldo yang telah dikurangi sebelumnya
		wallet.Balance += float64(withdrawal.Amount)
		_ = u.walletRepository.Update(wallet)
		return fmt.Errorf("failed to create withdrawal: %v", err)
	}
	return nil
}



func (u *withdrawUsecase) HistoryTransaction(withdrawalID int) ([]*domain.Withdrawal, error) {
	return u.withdrawRepository.HistoryWithdrawal(withdrawalID)
}