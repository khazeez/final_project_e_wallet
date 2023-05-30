package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestTransferRepository_FindOne(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTransferRepository(db, nil)

	transferID := 1
	expectedTransfer := &domain.Transfer{
		ID: transferID,
		SenderId: domain.SenderWallet{
			ID: transferID,
			UserId: domain.UserSender{
				Sender_ID:             1,
				Sender_Name:           "madee",
				Sender_Email:          "widhiwww",
				Sender_Password:       "wwww",
				Sender_ProfilePicture: "www",
				IsDeleted:             false,
			},
			Balance: 0,
		},
		ReceiferId: domain.ReceiverWallet{
			ID: transferID,
			UserId: domain.UserReceiver{
				Receifer_ID:             2,
				Receifer_Name:           "www",
				Receifer_Email:          "www",
				Receifer_Password:       "www",
				Receifer_ProfilePicture: "www",
				IsDeleted:               false,
			},
			Balance: 0,
		},
		Amount:    100.0,
		Timestamp: time.Time{},
	}

	rows := sqlmock.NewRows([]string{"transfer_id", "receiver_wallet_id", "amount", "balance", "user_id", "name", "email", "password"}).
		AddRow(transferID, expectedTransfer.ReceiferId.ID, expectedTransfer.Amount, 0.0, 1, "John Doe", "john@example.com", "password")

	query := "SELECT t.transfer_id, t.receiver_wallet_id, t.amount, s.balance, u.user_id, u.name, u.email, u.password FROM transfer t JOIN wallet s ON s.wallet_id = t.receiver_wallet_id JOIN users u ON s.user_id = u.user_id WHERE t.transfer_id = $1"
	mock.ExpectQuery(query).
		WithArgs(transferID).
		WillReturnRows(rows)

	query2 := "SELECT t.transfer_id2, t.sender_wallet_id, w.balance, u.user_id, u.name, u.email, u.password FROM transfer t JOIN wallet w ON t.sender_wallet_id = w.wallet_id JOIN users u ON w.user_id = u.user_id WHERE transfer_id = $1"
	mock.ExpectQuery(query2).
		WithArgs(transferID).
		WillReturnRows(rows)

	transfer, err := repo.FindOne(transferID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTransfer, transfer)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}



func TestTransferRepository_FindAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTransferRepository(db, nil)

	expectedTransfers := []*domain.Transfer{
		{ID: 1, Amount: 100.0, SenderId: domain.SenderWallet{ID: 1}, ReceiferId: domain.ReceiverWallet{ID: 2}},
		{ID: 2, Amount: 200.0, SenderId: domain.SenderWallet{ID: 3}, ReceiferId: domain.ReceiverWallet{ID: 4}},
	}

	rows := sqlmock.NewRows([]string{"transfer_id", "sender_id", "receiver_id", "amount", "timestamp"}).
		AddRow(1, expectedTransfers[0].SenderId.ID, expectedTransfers[0].ReceiferId.ID, expectedTransfers[0].Amount, time.Now()).
		AddRow(2, expectedTransfers[1].SenderId.ID, expectedTransfers[1].ReceiferId.ID, expectedTransfers[1].Amount, time.Now())

	mock.ExpectQuery("SELECT transfer_id, sender_id, receiver_id, amount, timestamp FROM Transfer").
		WillReturnRows(rows)

	transfers, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Equal(t, expectedTransfers, transfers)
}

func TestTransferRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTransferRepository(db, nil)

	transfer := &domain.Transfer{
		ID:         1,
		Amount:     100.0,
		SenderId:   domain.SenderWallet{ID: 1},
		ReceiferId: domain.ReceiverWallet{ID: 2},
	}

	mock.ExpectExec("UPDATE Transfer SET sender_id = ?").
		WithArgs(transfer.SenderId, transfer.ReceiferId, transfer.Amount, transfer.Timestamp, transfer.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(transfer)

	assert.NoError(t, err)
}

func TestTransferRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTransferRepository(db, nil)

	transferID := 1

	mock.ExpectExec("DELETE FROM Transfer WHERE transfer_id = $1").
		WithArgs(transferID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(transferID)

	assert.NoError(t, err)
}

func TestTransferRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTransferRepository(db, nil)

	transfer := &domain.Transfer{
		ID:         1,
		Amount:     100.0,
		SenderId:   domain.SenderWallet{ID: 1},
		ReceiferId: domain.ReceiverWallet{ID: 2},
	}

	mock.ExpectQuery("SELECT balance FROM Wallet WHERE wallet_id = $1").
		WithArgs(transfer.SenderId.ID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(200.0))

	mock.ExpectQuery("SELECT balance FROM Wallet WHERE wallet_id = $1").
		WithArgs(transfer.ReceiferId.ID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(0.0))

	mock.ExpectExec("UPDATE Wallet SET balance = balance - $1").
		WithArgs(transfer.Amount, transfer.SenderId.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("UPDATE Wallet SET balance = balance + $1").
		WithArgs(transfer.Amount, transfer.ReceiferId.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("INSERT INTO transfer").
		WithArgs(transfer.ID, transfer.SenderId.ID, transfer.ReceiferId.ID, transfer.Amount, transfer.Timestamp).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(transfer)

	assert.NoError(t, err)
}

func TestTransferRepository_History(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTransferRepository(db, nil)

	walletID := 1

	expectedTransfers := []*domain.Transfer{
		{ID: 1, SenderId: domain.SenderWallet{ID: 1}, ReceiferId: domain.ReceiverWallet{ID: 2}, Amount: 100.0},
		{ID: 2, SenderId: domain.SenderWallet{ID: 3}, ReceiferId: domain.ReceiverWallet{ID: 1}, Amount: 200.0},
	}

	rows := sqlmock.NewRows([]string{"transfer_id", "sender_wallet_id", "receiver_wallet_id", "amount", "timestamp", "balance"}).
		AddRow(expectedTransfers[0].ID, expectedTransfers[0].SenderId.ID, expectedTransfers[0].ReceiferId.ID, expectedTransfers[0].Amount, time.Now(), 100.0).
		AddRow(expectedTransfers[1].ID, expectedTransfers[1].SenderId.ID, expectedTransfers[1].ReceiferId.ID, expectedTransfers[1].Amount, time.Now(), 200.0)

	mock.ExpectQuery("SELECT t.transfer_id, t.receiver_wallet_id, t.amount, s.balance, u.user_id, u.name, u.email, u.password FROM transfer t JOIN wallet s ON s.wallet_id = t.receiver_wallet_id JOIN users u ON s.user_id = u.user_id WHERE t.transfer_id = $1").
		WithArgs(walletID).
		WillReturnRows(rows)

	transfers, err := repo.History(walletID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTransfers, transfers)
}
