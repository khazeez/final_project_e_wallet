package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)

type PaymentRepository interface {
	Create(payment *domain.Payment) error
	FindOne(paymentID int) (*domain.Payment, error)
	Update(payment *domain.Payment) error
	Delete(paymentID int) error
	HistoryPayment(wallet_id int) ([]*domain.Payment, error)
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (p *paymentRepository) Create(payment *domain.Payment) error {
	query := "SELECT balance FROM wallet WHERE wallet_id = $1"
	row := p.db.QueryRow(query, payment.WalletId.ID)
	var balance float64
	err := row.Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("wallet not found")
		}
		return fmt.Errorf("failed to get wallet balance: %v", err)
	}

	if balance < float64(payment.Amount) {
		return fmt.Errorf("insufficient balance")
	}

	balance = 0
	updateQuery := "UPDATE Wallet SET balance = balance - $1 WHERE wallet_id = $2"
	_, err = p.db.Exec(updateQuery, payment.Amount, payment.WalletId.ID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}

	time := time.Now()
	insertQuery := "INSERT INTO payment (payment_id, wallet_id, amount, timestamp, payment_type, payment_details) VALUES ($1, $2, $3, $4, $5, $6);"
	_, err = p.db.Exec(insertQuery, payment.ID, payment.WalletId.ID, payment.Amount, time, payment.PaymentType, payment.PaymentDetail)
	if err != nil {
		return fmt.Errorf("failed to create payment: %v", err)
	}
	return nil
}

func (p *paymentRepository) FindOne(paymentID int) (*domain.Payment, error) {
	query := `
	SELECT t.payment_id, t.amount, t.payment_type, t.payment_details, t.timestamp, w.wallet_id, u.user_id, u.name, u.email, u.password, u.profile_picture, u.is_deleted, w.balance
	FROM payment t
	JOIN wallet w ON t.wallet_id = w.wallet_id
	JOIN users u ON w.user_id = u.user_id
	WHERE t.payment_id = $1
	`

	row := p.db.QueryRow(query, paymentID)
	payment := &domain.Payment{}
	wallet := &domain.Wallet{}
	user := &domain.User{}

	err := row.Scan(
		&payment.ID,
		&payment.Amount,
		&payment.PaymentType,
		&payment.PaymentDetail,
		&payment.Timestamp,
		&wallet.ID,
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.ProfilePicture,
		&user.IsDeleted,
		&wallet.Balance,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to find payment: %v", err)
	}
	wallet.UserId = *user

	payment.WalletId = *wallet
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
	query := "DELETE FROM Payment WHERE payment_id = $1"
	_, err := r.db.Exec(query, paymentID)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %v", err)
	}
	return nil
}

func (r *paymentRepository) HistoryPayment(wallet_id int) ([]*domain.Payment, error) {
	query := `SELECT t.payment_id, t.amount, t.payment_type, t.payment_details, w.wallet_id, u.user_id, u.name, u.email, u.password, u.profile_picture, u.is_deleted, w.balance FROM payment t
	 	JOIN Wallet w ON t.wallet_id = w.wallet_id
	 	JOIN users u ON w.user_id = u.user_id
	 	WHERE t.wallet_id = $1`

	rows, err := r.db.Query(query, wallet_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %v", err)
	}
	defer rows.Close()

	payments := []*domain.Payment{}
	wallet := &domain.Wallet{}
	user := &domain.User{}

	for rows.Next() {
		payment := &domain.Payment{}
		err := rows.Scan(
			&payment.ID,
			&payment.Amount,
			&payment.PaymentType,
			&payment.PaymentDetail,
			&wallet.ID,
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.ProfilePicture,
			&user.IsDeleted,
			&wallet.Balance,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment row: %v", err)
		}

		wallet.UserId = *user
		payment.WalletId = *wallet

		payments = append(payments, payment)
	}
	return payments, nil
}
