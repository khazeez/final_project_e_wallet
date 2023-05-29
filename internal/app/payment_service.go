package app

import (
	// "fmt"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
)

type PaymentUsecase interface {
	CreatePayment(payment *domain.Payment) error
	GetPaymentByID(paymentID int) (*domain.Payment, error)
	UpdatePayment(payment *domain.Payment) error
	DeletePayment(paymentID int) error
	HistoryTransaction(paymentID int) ([]*domain.Payment, error)
	// MakePayment(payment *domain.Payment) error
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

func (u *paymentUsecase) HistoryTransaction(paymentID int) ([]*domain.Payment, error) {
	return u.paymentRepository.HistoryPayment(paymentID)
}

