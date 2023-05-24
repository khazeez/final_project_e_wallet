package app

import (
	"fmt"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
)

type PaymentUsecase interface {
	CreatePayment(payment *domain.Payment) error
	GetPaymentByID(paymentID int) (*domain.Payment, error)
	UpdatePayment(payment *domain.Payment) error
	DeletePayment(paymentID int) error
	MakePayment(payment *domain.Payment) error
}

type paymentUsecase struct {
	paymentRepository repository.PaymentRepository
	walletRepository  repository.WalletRepository
}

func NewPaymentUsecase(paymentRepository repository.PaymentRepository, walletRepository repository.WalletRepository) PaymentUsecase {
	return &paymentUsecase{
		paymentRepository: paymentRepository,
		walletRepository:  walletRepository,
	}
}

func (u *paymentUsecase) CreatePayment(payment *domain.Payment) error {
	return u.paymentRepository.Create(payment)
}

func (u *paymentUsecase) GetPaymentByID(paymentID int) (*domain.Payment, error) {
	return u.paymentRepository.FindOne(paymentID)
}

func (u *paymentUsecase) UpdatePayment(payment *domain.Payment) error {
	return u.paymentRepository.Update(payment)
}

func (u *paymentUsecase) DeletePayment(paymentID int) error {
	return u.paymentRepository.Delete(paymentID)
}

func (u *paymentUsecase) MakePayment(payment *domain.Payment) error {
	wallet, err := u.walletRepository.FindOne(payment.WalletId)
	if err != nil {
		return fmt.Errorf("failed to find wallet: %v", err)
	}

	if wallet.Balance < float64(payment.Amount) {
		return fmt.Errorf("insufficient balance in wallet")
	}

	// Kurangi saldo pada wallet
	wallet.Balance -= payment.Amount
	err = u.walletRepository.Update(wallet)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}

	// Simpan pembayaran
	err = u.paymentRepository.Create(payment)
	if err != nil {
		// Jika gagal menyimpan pembayaran, tambahkan kembali saldo yang telah dikurangi sebelumnya
		wallet.Balance += payment.Amount
		_ = u.walletRepository.Update(wallet)
		return fmt.Errorf("failed to create payment: %v", err)
	}

	return nil
}
