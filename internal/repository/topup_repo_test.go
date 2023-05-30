package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateTopup(t *testing.T) {
	// Membuat database mock dan mock controller
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat topupRepository dengan database mock
	topupRepo := NewTopupRepository(db)

	// Membuat data top-up untuk pengujian
	wallet := &domain.Wallet{
		ID:      1,
		Balance: 200.0,
		UserId: domain.User{
			ID:             1,
			Name:           "John Doe",
			Email:          "johndoe@example.com",
			Password:       "password",
			ProfilePicture: "profile.jpg",
			IsDeleted:      false,
		},
	}
	topup := &domain.TopUp{
		ID:        1,
		Amount:    100.0,
		Timestamp: time.Now(),
		WalletId:  *wallet,
	}

	// Expect query update saldo
	mock.ExpectExec(`UPDATE Wallet SET balance = balance \+ \$1 WHERE wallet_id = \$2`).
		WithArgs(topup.Amount, topup.WalletId.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Expect query insert top-up
	mock.ExpectExec(`INSERT INTO TopUp \(topup_id, wallet_id, amount, timestamp\) VALUES \(\$1, \$2, \$3, \$4\)`).
		WithArgs(topup.ID, topup.WalletId.ID, topup.Amount, topup.Timestamp).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Memanggil fungsi Create
	err = topupRepo.Create(topup)

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan semua expectation terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetLastTopupAmount(t *testing.T) {
	// Membuat database mock dan mock controller
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat topupRepository dengan database mock
	topupRepo := NewTopupRepository(db)

	// Membuat data walletID dan expectedAmount untuk pengujian
	walletID := 1
	expectedAmount := 100.0

	// Expect query select last top-up amount
	mock.ExpectQuery(`SELECT amount FROM TopUp WHERE wallet_id = \$1 ORDER BY topup_id DESC LIMIT 1`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(expectedAmount))

	// Memanggil fungsi GetLastTopupAmount
	amount, err := topupRepo.GetLastTopupAmount(walletID)

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan jumlah top-up terakhir sesuai dengan yang diharapkan
	assert.Equal(t, expectedAmount, amount)

	// Memastikan semua expectation terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindOneTopup(t *testing.T) {
	// Membuat database mock dan mock controller
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat topupRepository dengan database mock
	topupRepo := NewTopupRepository(db)

	// Membuat data topupID untuk pengujian
	topupID := 1

	// Membuat data top-up yang diharapkan
	expectedTopup := &domain.TopUp{
		ID:        1,
		WalletId:  domain.Wallet{},
		Amount:    100.0,
		Timestamp: time.Now(),
	}

	// Expect query select top-up
	mock.ExpectQuery(`SELECT t.topup_id, t.amount, t.timestamp, w.wallet_id, u.user_id, u.name, u.email, u.password, u.profile_picture, u.is_deleted, w.balance FROM TopUp t JOIN Wallet w ON t.wallet_id = w.wallet_id JOIN users u ON w.user_id = u.user_id WHERE t.topup_id = \$1`).
		WithArgs(topupID).
		WillReturnRows(sqlmock.NewRows([]string{"topup_id", "amount", "timestamp", "wallet_id", "user_id", "name", "email", "password", "profile_picture", "is_deleted", "balance"}).
			AddRow(
				expectedTopup.ID,
				expectedTopup.Amount,
				expectedTopup.Timestamp,
				expectedTopup.WalletId.ID,
				expectedTopup.WalletId.UserId.ID,
				expectedTopup.WalletId.UserId.Name,
				expectedTopup.WalletId.UserId.Email,
				expectedTopup.WalletId.UserId.Password,
				expectedTopup.WalletId.UserId.ProfilePicture,
				expectedTopup.WalletId.UserId.IsDeleted,
				expectedTopup.WalletId.Balance,
			))

	// Memanggil fungsi FindOne
	result, err := topupRepo.FindOne(topupID)

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan data top-up yang ditemukan sesuai dengan yang diharapkan
	assert.Equal(t, expectedTopup, result)

	// Memastikan semua expectation terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTopup(t *testing.T) {
	// Membuat database mock dan mock controller
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat topupRepository dengan database mock
	topupRepo := NewTopupRepository(db)

	// Membuat data top-up untuk pengujian
	topup := &domain.TopUp{
		ID:     1,
		Amount: 100.0,
	}

	// Expect query select previous top-up amount
	mock.ExpectQuery(`SELECT amount FROM topup WHERE topup_id = \$1`).
		WithArgs(topup.ID).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(topup.Amount))

	// Expect query update top-up amount
	mock.ExpectExec(`UPDATE topup SET amount = \$1 WHERE topup_id = \$2`).
		WithArgs(topup.Amount+topup.Amount, topup.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Memanggil fungsi Update
	err = topupRepo.Update(topup)

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan semua expectation terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTopup(t *testing.T) {
	// Membuat database mock dan mock controller
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat topupRepository dengan database mock
	topupRepo := NewTopupRepository(db)

	// Membuat data topupID untuk pengujian
	topupID := 1

	// Expect query delete top-up
	mock.ExpectExec(`DELETE FROM TopUp WHERE topup_id = \$1`).
		WithArgs(topupID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Memanggil fungsi Delete
	err = topupRepo.Delete(topupID)

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan semua expectation terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHistoryTopup(t *testing.T) {
	// Membuat database mock dan mock controller
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat topupRepository dengan database mock
	topupRepo := NewTopupRepository(db)

	// Membuat data walletID untuk pengujian
	walletID := 1

	// Membuat data top-up yang diharapkan
	expectedTopups := []*domain.TopUp{
		{
			ID:        1,
		WalletId:  domain.Wallet{},
		Amount:    100.0,
		Timestamp: time.Time{},
		},
		{
			ID:        2,
			WalletId:  domain.Wallet{},
			Amount:    100.0,
			Timestamp: time.Time{},
		},
	}

	// Expect query select history top-ups
	mock.ExpectQuery(`SELECT t.topup_id, t.amount, w.wallet_id, u.user_id, u.name, u.email, u.password, u.profile_picture, u.is_deleted, w.balance FROM topup t JOIN Wallet w ON t.wallet_id = w.wallet_id JOIN users u ON w.user_id = u.user_id WHERE t.wallet_id = \$1`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"topup_id", "amount", "wallet_id", "user_id", "name", "email", "password", "profile_picture", "is_deleted", "balance"}).
			AddRow(
				expectedTopups[0].ID,
				expectedTopups[0].Amount,
				expectedTopups[0].WalletId.ID,
				expectedTopups[0].WalletId.UserId.ID,
				expectedTopups[0].WalletId.UserId.Name,
				expectedTopups[0].WalletId.UserId.Email,
				expectedTopups[0].WalletId.UserId.Password,
				expectedTopups[0].WalletId.UserId.ProfilePicture,
				expectedTopups[0].WalletId.UserId.IsDeleted,
				expectedTopups[0].WalletId.Balance,
			).
			AddRow(
				expectedTopups[1].ID,
				expectedTopups[1].Amount,
				expectedTopups[1].WalletId.ID,
				expectedTopups[1].WalletId.UserId.ID,
				expectedTopups[1].WalletId.UserId.Name,
				expectedTopups[1].WalletId.UserId.Email,
				expectedTopups[1].WalletId.UserId.Password,
				expectedTopups[1].WalletId.UserId.ProfilePicture,
				expectedTopups[1].WalletId.UserId.IsDeleted,
				expectedTopups[1].WalletId.Balance,
			))

	// Memanggil fungsi HistoryTopup
	result, err := topupRepo.HistoryTopup(walletID)

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan data top-up yang ditemukan sesuai dengan yang diharapkan
	assert.Equal(t, expectedTopups, result)

	// Memastikan semua expectation terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())
}
