package repository

import (
	"database/sql"
	"fmt"

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
	query := "SELECT t.transfer_id, t.receiver_wallet_id, t.amount, s.balance, u.user_id, u.name, u.email, u.password FROM transfer t JOIN wallet s ON s.wallet_id = t.receiver_wallet_id JOIN users u ON s.user_id = u.user_id WHERE t.transfer_id = $1;"

	query2 := "SELECT t.transfer_id, t.sender_wallet_id, w.balance, u.user_id, u.name, u.email, u.password FROM transfer t JOIN wallet w ON t.sender_wallet_id = w.wallet_id JOIN users u ON w.user_id = u.user_id WHERE transfer_id = $1"

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

		&transfer.ID,
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
	// Cek apakah wallet dengan ID yang diberikan ada dalam database
	query := "SELECT balance FROM Wallet WHERE wallet_id = $1"
	row := r.db.QueryRow(query, transfer.SenderId.ID)
	var balance domain.Transfer
	err := row.Scan(&balance.SenderId.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("id not found")
		}
		return fmt.Errorf("failed to get wallet balance AAAAAAAAAA: %v", err)
	}

	query2 := "SELECT balance FROM Wallet WHERE wallet_id = $1"
	row2 := r.db.QueryRow(query2, transfer.ReceiferId.ID)
	var balance2 domain.Transfer
	err = row2.Scan(&balance2.ReceiferId.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("id not found")
		}
		return fmt.Errorf("failed to get wallet balance BBBBBBBBB: %v", err)
	}

	if balance.SenderId.Balance < float64(transfer.Amount) {
		return fmt.Errorf("insufficient balance")
	}

	balance.SenderId.Balance = 0
	updateQuerySender := "UPDATE Wallet SET balance = balance - $1 WHERE wallet_id = $2"
	_, err = r.db.Exec(updateQuerySender, transfer.Amount, transfer.SenderId.ID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}
	balance.ReceiferId.Balance = 0
	updateQueryReceifer := "UPDATE Wallet SET balance = balance + $1 WHERE wallet_id = $2"
	_, err = r.db.Exec(updateQueryReceifer, transfer.Amount, transfer.ReceiferId.ID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}

	insertQuery := "INSERT INTO transfer (transfer_id, sender_wallet_id, receiver_wallet_id, amount, timestamp) VALUES ($1, $2, $3, $4, $5)"
	_, err = r.db.Exec(insertQuery, transfer.ID, transfer.SenderId.ID, transfer.ReceiferId.ID, transfer.Amount, transfer.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to transfer: %v", err)
	}
	return nil
}
func (r *transferRepository) History(walletID int) ([]*domain.Transfer, error) {

	query := `SELECT t.transfer_id, t.sender_wallet_id, t.receiver_wallet_id, t.amount, t.timestamp, w.balance
	FROM transfer t
	JOIN wallet w ON t.sender_wallet_id = w.wallet_id
	JOIN users u ON w.user_id = u.user_id
	WHERE t.sender_wallet_id = $1 OR t.receiver_wallet_id = $1`

	rows, err := r.db.Query(query, walletID)
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

