package repository

import (
	"database/sql"
	"fmt"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)

type WalletRepository interface {
	Create(wallet *domain.Wallet) error
	FindOne(walletID int) (*domain.Wallet, error)
	Update(wallet *domain.Wallet) error
	Delete(walletID int) error

}
type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (r *walletRepository) Create(wallet *domain.Wallet) error {
	query := "INSERT INTO Wallet (wallet_id, user_id, balance) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, wallet.ID, wallet.UserId, wallet.Balance)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %v", err)
	}
	return nil
}


func (r *walletRepository) FindOne(walletID int) (*domain.Wallet, error) {
	query := "SELECT wallet_id, user_id, balance FROM Wallet WHERE wallet_id = ?"
	row := r.db.QueryRow(query, walletID)
	wallet := &domain.Wallet{}
	err := row.Scan(&wallet.ID, &wallet.UserId, &wallet.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Dompet tidak ditemukan
		}
		return nil, fmt.Errorf("failed to find wallet: %v", err)
	}
	return wallet, nil
}



func (r *walletRepository) Update(wallet *domain.Wallet) error {
	// Mengambil saldo top-up terakhir untuk wallet ini
	queryTopUp := "SELECT amount FROM TopUp WHERE wallet_id = ? ORDER BY timestamp DESC LIMIT 1"
	row := r.db.QueryRow(queryTopUp, wallet.ID)
	var topupAmount float64
	err := row.Scan(&topupAmount)
	if err != nil {
		return fmt.Errorf("failed to get top-up amount: %v", err)
	}

	// Menambahkan saldo top-up ke balance wallet
	wallet.Balance += topupAmount

	// Perbarui nilai balance pada tabel wallet dengan nilai baru
	queryUpdate := "UPDATE Wallet SET balance = ? WHERE wallet_id = ?"
	_, err = r.db.Exec(queryUpdate, wallet.Balance, wallet.ID)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %v", err)
	}

	return nil
}

func (r *walletRepository) Delete(walletID int) error {
	query := "DELETE FROM Wallet WHERE wallet_id = ?"
	_, err := r.db.Exec(query, walletID)
	if err != nil {
		return fmt.Errorf("failed to delete wallet: %v", err)
	}
	return nil
}
