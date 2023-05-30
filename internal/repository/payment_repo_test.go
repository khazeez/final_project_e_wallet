package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestPaymentRepository_Create(t *testing.T) { //belum ok
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewPaymentRepository(db)

	payment := &domain.Payment{
		ID:           1,
		WalletId:    domain.Wallet{ID: 1},
		Amount:       100.0,
		Timestamp:    time.Now(),
		PaymentType:  "Credit",
		PaymentDetail: "Payment for a product",
	}

	mock.ExpectQuery("SELECT balance FROM wallet").
		WithArgs(payment.WalletId.ID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(200.0))

	mock.ExpectExec("UPDATE Wallet SET balance").
		WithArgs(payment.Amount, payment.WalletId.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("INSERT INTO payment").
		WithArgs(payment.ID, payment.WalletId.ID, payment.Amount, payment.Timestamp, payment.PaymentType, payment.PaymentDetail).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Create(payment)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentRepository_FindOne(t *testing.T) {  // OK
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewPaymentRepository(db)

	paymentID := 1

	rows := sqlmock.NewRows([]string{
		"payment_id", "amount", "payment_type", "payment_details", "timestamp",
		"wallet_id", "user_id", "name", "email", "password", "profile_picture", "is_deleted", "balance",
	}).AddRow(
		paymentID, 100.0, "Credit", "Payment for a product", time.Now(),
		1, 1, "John Doe", "john.doe@example.com", "password123", "profile.jpg", false, 200.0,
	)

	mock.ExpectQuery("SELECT t.payment_id, t.amount, t.payment_type, t.payment_details, t.timestamp, w.wallet_id, u.user_id, u.name, u.email, u.password, u.profile_picture, u.is_deleted, w.balance FROM payment").
		WithArgs(paymentID).
		WillReturnRows(rows)

	payment, err := repo.FindOne(paymentID)
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, paymentID, payment.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentRepository_Update(t *testing.T) { // belum ok
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewPaymentRepository(db)

	payment := &domain.Payment{
		ID:           1,
		WalletId:     domain.Wallet{ID: 1},
		Amount:       100.0,
		Timestamp:    time.Now(),
		PaymentType:  "Credit",
		PaymentDetail: "Payment for a product",
	}

	mock.ExpectExec("UPDATE Payment SET wallet_id").
		WithArgs(payment.WalletId, payment.Amount, payment.Timestamp, payment.PaymentType, payment.PaymentDetail, payment.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(payment)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentRepository_Delete(t *testing.T) { //OK
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewPaymentRepository(db)

	paymentID := 1

	mock.ExpectExec("DELETE FROM Payment").
		WithArgs(paymentID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(paymentID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentRepository_HistoryPayment(t *testing.T) { //OK
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewPaymentRepository(db)

	walletID := 1

	rows := sqlmock.NewRows([]string{
		"payment_id", "amount", "payment_type", "payment_details",
		"wallet_id", "user_id", "name", "email", "password", "profile_picture", "is_deleted", "balance",
	}).AddRow(
		1, 100.0, "Credit", "Payment for a product",
		walletID, 1, "John Doe", "john.doe@example.com", "password123", "profile.jpg", false, 200.0,
	)

	mock.ExpectQuery("SELECT t.payment_id, t.amount, t.payment_type, t.payment_details, w.wallet_id, u.user_id, u.name, u.email, u.password, u.profile_picture, u.is_deleted, w.balance FROM payment").
		WithArgs(walletID).
		WillReturnRows(rows)

	payments, err := repo.HistoryPayment(walletID)
	assert.NoError(t, err)
	assert.NotNil(t, payments)
	assert.Len(t, payments, 1)

	assert.NoError(t, mock.ExpectationsWereMet())
}