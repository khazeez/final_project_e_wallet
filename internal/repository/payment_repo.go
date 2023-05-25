package repository

import (
	"database/sql"
	"fmt"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)




type PaymentRepository interface {
	Create(payment *domain.Payment) error
	FindOne(paymentID int) (*domain.Payment, error)
	Update(payment *domain.Payment) error
	Delete(paymentID int) error
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	query := "INSERT INTO Payment (payment_id, wallet_id, amount, timestamp, payment_type, payment_detail) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, payment.ID, payment.WalletId, payment.Amount, payment.Timestamp, payment.PaymentType, payment.PaymentDetail)
	if err != nil {
		return fmt.Errorf("failed to create payment: %v", err)
	}
	return nil
}

func (r *paymentRepository) FindOne(paymentID int) (*domain.Payment, error) {
	query := "SELECT payment_id, wallet_id, amount, timestamp, payment_type, payment_detail FROM Payment WHERE payment_id = ?"
	row := r.db.QueryRow(query, paymentID)
	payment := &domain.Payment{}
	err := row.Scan(&payment.ID, &payment.WalletId, &payment.Amount, &payment.Timestamp, &payment.PaymentType, &payment.PaymentDetail)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to find payment: %v", err)
	}
	return payment, nil
}

func (r *paymentRepository) Update(payment *domain.Payment) error {
	query := "UPDATE Payment SET wallet_id = ?, amount = ?, timestamp = ?, payment_type = ?, payment_detail = ? WHERE payment_id = ?"
	_, err := r.db.Exec(query, payment.WalletId, payment.Amount, payment.Timestamp, payment.PaymentType, payment.PaymentDetail, payment.ID)
	if err != nil {
		return fmt.Errorf("failed to update payment: %v", err)
	}
	return nil
}

func (r *paymentRepository) Delete(paymentID int) error {
	query := "DELETE FROM Payment WHERE payment_id = ?"
	_, err := r.db.Exec(query, paymentID)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %v", err)
	}
	return nil
}
