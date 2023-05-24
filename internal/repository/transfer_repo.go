package repository

import (
	"database/sql"
	"fmt"
	"errors"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)

type TransferRepository interface {
	Create(transfer *domain.Transfer) error
	FindOne(transferID int) (*domain.Transfer, error)
	FindAll() ([]*domain.Transfer, error)
	Update(transfer *domain.Transfer) error
	Delete(transferID int) error
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

func (r *transferRepository) Create(transfer *domain.Transfer) error {
	senderWallet, err := r.walletRepository.FindOne(transfer.SenderId.ID)
	if err != nil {
		return fmt.Errorf("failed to get sender wallet: %v", err)
	}

	receiverWallet, err := r.walletRepository.FindOne(transfer.ReceiferId.ID)
	if err != nil {
		return fmt.Errorf("failed to get receiver wallet: %v", err)
	}

	// Validasi saldo cukup pada wallet pengirim
	if senderWallet.Balance < float64(transfer.Amount) {
		return errors.New("insufficient balance for transfer")
	}

	// Mengurangi saldo pada wallet pengirim
	senderWallet.Balance -= transfer.Amount

	// Menambah saldo pada wallet penerima
	receiverWallet.Balance += transfer.Amount

	// Memperbarui saldo pada tabel Wallet
	err = r.walletRepository.Update(senderWallet)
	if err != nil {
		return fmt.Errorf("failed to update sender wallet: %v", err)
	}

	err = r.walletRepository.Update(receiverWallet)
	if err != nil {
		return fmt.Errorf("failed to update receiver wallet: %v", err)
	}

	// Melakukan penyimpanan transfer ke tabel Transfer
	query := "INSERT INTO Transfer (sender_id, receiver_id, amount, timestamp) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(query, transfer.SenderId, transfer.ReceiferId, transfer.Amount, transfer.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to create transfer: %v", err)
	}

	transferID, _ := result.LastInsertId()
	transfer.ID = int(transferID)
	return nil
}


func (r *transferRepository) FindOne(transferID int) (*domain.Transfer, error) {
	query := "SELECT transfer_id, sender_id, receiver_id, amount, timestamp FROM Transfer WHERE transfer_id = ?"
	row := r.db.QueryRow(query, transferID)
	transfer := &domain.Transfer{}
	err := row.Scan(&transfer.ID, &transfer.SenderId, &transfer.ReceiferId, &transfer.Amount, &transfer.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transfer not found")
		}
		return nil, fmt.Errorf("failed to get transfer: %v", err)
	}
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
	query := "DELETE FROM Transfer WHERE transfer_id = ?"
	_, err := r.db.Exec(query, transferID)
	if err != nil {
		return fmt.Errorf("failed to delete transfer: %v", err)
	}
	return nil
}
