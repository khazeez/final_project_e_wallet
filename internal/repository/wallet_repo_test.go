// package repository

// import (

// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
// )

// func TestWalletRepository_Create(t *testing.T) {
// 	// Membuat database mock dan mock controller
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// Membuat walletRepository dengan database mock
// 	walletRepo := NewWalletRepository(db)

// 	// Membuat data wallet yang akan di-create
// 	wallet := &domain.Wallet{
// 		ID:      1,
// 		UserId:  domain.User{ID: 1},
// 		Balance: 100.0,
// 	}

// 	// Expect query insert wallet
// 	mock.ExpectExec(`INSERT INTO Wallet (.+) VALUES (.+)`).
// 		WithArgs(wallet.ID, wallet.UserId.ID, wallet.Balance).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	// Memanggil fungsi Create
// 	err = walletRepo.Create(wallet)

// 	// Memastikan tidak ada error
// 	assert.NoError(t, err)

// 	// Memastikan semua expectation terpenuhi
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestWalletRepository_FindOne(t *testing.T) {
// 	// Membuat database mock dan mock controller
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// Membuat walletRepository dengan database mock
// 	walletRepo := NewWalletRepository(db)

// 	// Membuat data walletID untuk pengujian
// 	walletID := 1

// 	// Membuat data wallet yang diharapkan
// 	expectedWallet := &domain.Wallet{
// 		ID:      1,
// 		UserId:  domain.User{ID: 1, Name: "John Doe", Email: "john@example.com"},
// 		Balance: 100.0,
// 	}

// 	// Expect query select wallet
// 	mock.ExpectQuery(`SELECT (.+) FROM Wallet (.+) WHERE (.+)`).
// 		WithArgs(walletID).
// 		WillReturnRows(sqlmock.NewRows([]string{"wallet_id", "balance", "user_id", "name", "email"}).
// 			AddRow(
// 				expectedWallet.ID,
// 				expectedWallet.Balance,
// 				expectedWallet.UserId.ID,
// 				expectedWallet.UserId.Name,
// 				expectedWallet.UserId.Email,
// 			))

// 	// Memanggil fungsi FindOne
// 	result, err := walletRepo.FindOne(walletID)

// 	// Memastikan tidak ada error
// 	assert.NoError(t, err)

// 	// Memastikan data wallet yang ditemukan sesuai dengan yang diharapkan
// 	assert.Equal(t, expectedWallet, result)

// 	// Memastikan semua expectation terpenuhi
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }


// func TestWalletRepository_Delete(t *testing.T) {
// 	// Membuat database mock dan mock controller
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// Membuat walletRepository dengan database mock
// 	walletRepo := NewWalletRepository(db)

// 	// Membuat data walletID untuk pengujian
// 	walletID := 1

// 	// Expect query delete wallet
// 	mock.ExpectExec(`DELETE FROM Wallet WHERE (.+)`).
// 		WithArgs(walletID).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	// Memanggil fungsi Delete
// 	err = walletRepo.Delete(walletID)

// 	// Memastikan tidak ada error
// 	assert.NoError(t, err)

// 	// Memastikan semua expectation terpenuhi
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

package repository_test

import (
	// "database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository dengan db mock
	repo := repository.NewWalletRepository(db)

	wallet := &domain.Wallet{
		ID:      1,
		UserId:  domain.User{ID: 1},
		Balance: 1000,
	}

	// Mengharapkan eksekusi query INSERT yang berhasil
	mock.ExpectExec("INSERT INTO Wallet").WithArgs(wallet.ID, wallet.UserId.ID, wallet.Balance).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(wallet)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository dengan db mock
	repo := repository.NewWalletRepository(db)

	wallet := &domain.Wallet{
		ID:      1,
		UserId:  domain.User{ID: 1},
		Balance: 1000,
	}

	// Mengharapkan error saat eksekusi query INSERT gagal
	mock.ExpectExec("INSERT INTO Wallet").WithArgs(wallet.ID, wallet.UserId.ID, wallet.Balance).WillReturnError(errors.New("database error"))

	err = repo.Create(wallet)
	assert.EqualError(t, err, "failed to create wallet: database error")

	assert.NoError(t, mock.ExpectationsWereMet())
}


func TestFindOne_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository dengan db mock
	repo := repository.NewWalletRepository(db)

	walletID := 1

	// Mengharapkan eksekusi query SELECT yang tidak mengembalikan data wallet
	mock.ExpectQuery("SELECT (.+) FROM Wallet").WithArgs(walletID).WillReturnRows(sqlmock.NewRows([]string{}))

	wallet, err := repo.FindOne(walletID)
	assert.NoError(t, err)
	assert.Nil(t, wallet)

	assert.NoError(t, mock.ExpectationsWereMet())
}


func TestUpdate_GetTopUpAmountError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository dengan db mock
	repo := repository.NewWalletRepository(db)

	walletID := 1

	// Mengharapkan error saat eksekusi query SELECT untuk mendapatkan top-up amount gagal
	mock.ExpectQuery("SELECT amount FROM TopUp").WithArgs(walletID).WillReturnError(errors.New("database error"))

	wallet := &domain.Wallet{
		ID:      walletID,
		UserId:  domain.User{ID: 1},
		Balance: 1000,
	}

	err = repo.Update(wallet)
	assert.EqualError(t, err, "failed to get top-up amount: database error")

	assert.NoError(t, mock.ExpectationsWereMet())
}


func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository dengan db mock
	repo := repository.NewWalletRepository(db)

	walletID := 1

	// Mengharapkan eksekusi query DELETE yang berhasil
	mock.ExpectExec("DELETE FROM Wallet").WithArgs(walletID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(walletID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository dengan db mock
	repo := repository.NewWalletRepository(db)

	walletID := 1

	// Mengharapkan error saat eksekusi query DELETE gagal
	mock.ExpectExec("DELETE FROM Wallet").WithArgs(walletID).WillReturnError(errors.New("database error"))

	err = repo.Delete(walletID)
	assert.EqualError(t, err, "failed to delete wallet: database error")

	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestWalletRepository_Update(t *testing.T) {
	// Membuat database mock dan mock controller
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat walletRepository dengan database mock
	walletRepo := repository.NewWalletRepository(db)

	// Membuat data wallet untuk pengujian
	wallet := &domain.Wallet{
		ID:      1,
		UserId:  domain.User{ID: 1},
		Balance: 100.0,
	}

	// Membuat data saldo top-up yang diharapkan
	expectedTopupAmount := 50.0

	// Expect query select top-up
	mock.ExpectQuery(`SELECT amount FROM TopUp WHERE wallet_id = \$1 ORDER BY timestamp DESC LIMIT 1`).
	WithArgs(wallet.ID).
	WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(expectedTopupAmount))

	// Expect query update wallet
	mock.ExpectExec(`UPDATE Wallet SET (.+) WHERE (.+)`).
		WithArgs(wallet.Balance+expectedTopupAmount, wallet.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Memanggil fungsi Update
	err = walletRepo.Update(wallet)

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan semua expectation terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWalletRepository_FindOne(t *testing.T) {
	// Membuat database mock dan mock controller
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat walletRepository dengan database mock
	walletRepo := repository.NewWalletRepository(db)

	// Membuat data walletID untuk pengujian
	walletID := 1

	// Membuat data wallet yang diharapkan
	expectedWallet := &domain.Wallet{
		ID:      1,
		UserId:  domain.User{ID: 1, Name: "John Doe", Email: "john@example.com"},
		Balance: 100.0,
	}

	// Expect query select wallet
	mock.ExpectQuery(`SELECT (.+) FROM Wallet (.+) JOIN(.+) WHERE (.+)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"wallet_id", "balance", "user_id", "name", "email","password","profile_picture","is_delete"}).
			AddRow(
				expectedWallet.ID,
				expectedWallet.Balance,
				expectedWallet.UserId.ID,
				expectedWallet.UserId.Name,
				expectedWallet.UserId.Email,
				expectedWallet.UserId.Password,
				expectedWallet.UserId.ProfilePicture,
				expectedWallet.UserId.IsDeleted,

			))

	// Memanggil fungsi FindOne
	result, err := walletRepo.FindOne(walletID)

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan data wallet yang ditemukan sesuai dengan yang diharapkan
	assert.NotNil(t, result)
	assert.Equal(t, expectedWallet.ID, result.ID)
	assert.Equal(t, expectedWallet.UserId.ID, result.UserId.ID)
	assert.Equal(t, expectedWallet.UserId.Name, result.UserId.Name)
	assert.Equal(t, expectedWallet.UserId.Email, result.UserId.Email)
	assert.Equal(t, expectedWallet.Balance, result.Balance)

	// Memastikan semua expectation terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())
}


