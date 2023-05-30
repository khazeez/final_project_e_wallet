package repository

import (
	"database/sql"
	"fmt"
	"time"

	// "errors"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)

type TransferRepository interface {
	Create(transfer *domain.Transfer) error
	FindOne(transferID int) (*domain.Transfer, error)
	FindAll() ([]*domain.Transfer, error)
	Update(transfer *domain.Transfer) error
	Delete(transferID int) error
	History(wallet_id int) ([]*domain.Transfer, error)

}
type transferRepository struct {
	db               *sql.DB
	walletRepository WalletRepository // Menggunakan WalletRepository yang sudah didefinisikan
}

func NewTransferRepository(db *sql.DB, walletRepository WalletRepository) TransferRepository {
	return &transferRepository{
		db:               db,
		walletRepository: walletRepository,
	}
}

func (r *transferRepository) FindOne(transferID int) (*domain.Transfer, error) {
	query := "SELECT t.transfer_id, t.receiver_wallet_id, t.amount, s.balance, u.user_id, u.name, u.email, u.password FROM transfer t JOIN wallet s ON s.wallet_id = t.receiver_wallet_id JOIN users u ON s.user_id = u.user_id WHERE t.transfer_id = $1"

	query2 := "SELECT t.transfer_id2, t.sender_wallet_id, w.balance, u.user_id, u.name, u.email, u.password FROM transfer t JOIN wallet w ON t.sender_wallet_id = w.wallet_id JOIN users u ON w.user_id = u.user_id WHERE transfer_id = $1"

	row := r.db.QueryRow(query, transferID)
	row2 := r.db.QueryRow(query2, transferID)

	transfer := &domain.Transfer{}
	walletSender := &domain.SenderWallet{}
	walletReceifer := &domain.ReceiverWallet{}
	sender := &domain.UserSender{}
	receifer := &domain.UserReceiver{}

	err := row.Scan(
		&transfer.ID,
		&walletReceifer.ID,
		&transfer.Amount,
		&walletReceifer.Balance,
		&receifer.Receifer_ID,
		&receifer.Receifer_Name,
		&receifer.Receifer_Email,
		&receifer.Receifer_Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transfer not found")
		}
		return nil, fmt.Errorf("failed to get transfer: %v", err)
	}

	err2 := row2.Scan(

		&transfer.SenderId.ID,
		&walletSender.ID,
		&walletSender.Balance,
		&sender.Sender_ID,
		&sender.Sender_Name,
		&sender.Sender_Email,
		&sender.Sender_Password)
	if err2 != nil {
		panic(err2)
	}

	walletSender.UserId = *sender
	transfer.SenderId = *walletSender
	walletReceifer.UserId = *receifer
	transfer.ReceiferId = *walletReceifer
	return transfer, nil
}

func (r *transferRepository) FindAll() ([]*domain.Transfer, error) {
	query := "SELECT transfer_id, sender_id, receiver_id, amount, timestamp FROM Transfer"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfers: %v", err)
	}
	defer rows.Close()

	transfers := []*domain.Transfer{}
	for rows.Next() {
		transfer := &domain.Transfer{}
		err := rows.Scan(&transfer.ID, &transfer.SenderId, &transfer.ReceiferId, &transfer.Amount, &transfer.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transfer row: %v", err)
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func (r *transferRepository) Update(transfer *domain.Transfer) error {
	query := "UPDATE Transfer SET sender_id = ?, receiver_id = ?, amount = ?, timestamp = ? WHERE transfer_id = ?"
	_, err := r.db.Exec(query, transfer.SenderId, transfer.ReceiferId, transfer.Amount, transfer.Timestamp, transfer.ID)
	if err != nil {
		return fmt.Errorf("failed to update transfer: %v", err)
	}
	return nil
}

func (r *transferRepository) Delete(transferID int) error {
	query := "DELETE FROM Transfer WHERE transfer_id = $1"
	_, err := r.db.Exec(query, transferID)
	if err != nil {
		return fmt.Errorf("failed to delete transfer: %v", err)
	}
	return nil
}

func (r *transferRepository) Create(transfer *domain.Transfer) error {
	tx, err := r.db.Begin() // Mulai transaksi
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Cek apakah wallet pengirim dengan ID yang diberikan ada dalam database
	querySender := "SELECT balance FROM Wallet WHERE wallet_id = $1"
	rowSender := tx.QueryRow(querySender, transfer.SenderId.ID)
	var senderBalance float64
	err = rowSender.Scan(&senderBalance)
	if err != nil {
		_ = tx.Rollback() // Batalkan transaksi
		if err == sql.ErrNoRows {
			return fmt.Errorf("sender wallet not found")
		}
		return fmt.Errorf("failed to get sender wallet balance: %v", err)
	}

	// Cek apakah wallet penerima dengan ID yang diberikan ada dalam database
	queryReceiver := "SELECT balance FROM Wallet WHERE wallet_id = $1"
	rowReceiver := tx.QueryRow(queryReceiver, transfer.ReceiferId.ID)
	var receiverBalance float64
	err = rowReceiver.Scan(&receiverBalance)
	if err != nil {
		_ = tx.Rollback() // Batalkan transaksi
		if err == sql.ErrNoRows {
			return fmt.Errorf("receiver wallet not found")
		}
		return fmt.Errorf("failed to get receiver wallet balance: %v", err)
	}

	// Cek apakah saldo wallet pengirim cukup untuk melakukan transfer
	if senderBalance < float64(transfer.Amount) {
		_ = tx.Rollback() // Batalkan transaksi
		return fmt.Errorf("insufficient balance")
	}

	// Kurangi saldo wallet pengirim sesuai dengan jumlah transfer
	updateQuerySender := "UPDATE Wallet SET balance = balance - $1 WHERE wallet_id = $2"
	_, err = tx.Exec(updateQuerySender, transfer.Amount, transfer.SenderId.ID)
	if err != nil {
		_ = tx.Rollback() // Batalkan transaksi
		return fmt.Errorf("failed to update sender wallet balance: %v", err)
	}

	// Tambah saldo wallet penerima sesuai dengan jumlah transfer
	updateQueryReceiver := "UPDATE Wallet SET balance = balance + $1 WHERE wallet_id = $2"
	_, err = tx.Exec(updateQueryReceiver, transfer.Amount, transfer.ReceiferId.ID)
	if err != nil {
		_ = tx.Rollback() // Batalkan transaksi
		return fmt.Errorf("failed to update receiver wallet balance: %v", err)
	}

	// Simpan data transfer ke dalam tabel Transfer
	time:=time.Now()
	insertQuery := "INSERT INTO Transfer (transfer_id, sender_wallet_id, receiver_wallet_id, amount, timestamp) VALUES ($1, $2, $3, $4, $5)"
	_, err = tx.Exec(insertQuery, transfer.ID, transfer.SenderId.ID, transfer.ReceiferId.ID, transfer.Amount, time)
	if err != nil {
		_ = tx.Rollback() // Batalkan transaksi
		return fmt.Errorf("failed to create transfer: %v", err)
	}

	err = tx.Commit() // Konfirmasi transaksi
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *transferRepository) History(walletID int) ([]*domain.Transfer, error) {

	query1 := `SELECT t.transfer_id, t.sender_wallet_id, t.receiver_wallet_id, t.amount, t.timestamp, w.balance, u.name
	FROM transfer t
	JOIN wallet w ON t.sender_wallet_id = w.wallet_id
	JOIN users u ON w.user_id = u.user_id
	WHERE t.sender_wallet_id = $1 OR t.receiver_wallet_id = $1`
	query2 := `SELECT t.transfer_id, t.sender_wallet_id, t.receiver_wallet_id, t.amount, t.timestamp, w.balance, u.name
	FROM transfer t
	JOIN wallet w ON t.sender_wallet_id = w.wallet_id
	JOIN users u ON w.user_id = u.user_id
	WHERE t.sender_wallet_id = $1 OR t.receiver_wallet_id = $1`

	rows, err := r.db.Query(query1,query2, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfers: %v", err)
	}
	defer rows.Close()

	transfers := []*domain.Transfer{}

	for rows.Next() {
		transfer := &domain.Transfer{}

		senderWallet := &domain.SenderWallet{}
		receiverWallet := &domain.ReceiverWallet{}
		senderUser := &domain.UserSender{}
		receiverUser := &domain.UserReceiver{}

		err := rows.Scan(
			&transfer.ID,
			&senderWallet.ID,
			&receiverWallet.ID,
			&transfer.Amount,
			&transfer.Timestamp,
			&senderWallet.Balance,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transfer row: %v", err)
		}
		err = rows.Scan(
			&transfer.ID,
			&senderWallet.ID,
			&receiverWallet.ID,
			&transfer.Amount,
			&transfer.Timestamp,
			&senderWallet.Balance,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transfer row: %v", err)
		}

		senderWallet.UserId = *senderUser
		transfer.SenderId = *senderWallet

		if receiverWallet.ID != 0 {
			receiverUser.Receifer_ID = receiverWallet.UserId.Receifer_ID
			receiverWallet.UserId = *receiverUser
			transfer.ReceiferId = *receiverWallet
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

