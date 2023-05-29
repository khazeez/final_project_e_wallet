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
	query := "INSERT INTO Wallet (wallet_id, user_id, balance) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, wallet.ID, wallet.UserId.ID, wallet.Balance)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %v", err)
	}
	return nil
}



func (r *walletRepository) FindOne(walletID int) (*domain.Wallet, error) {
	query := `
		SELECT w.wallet_id, w.balance, u.user_id, u.name, u.email, u.password, u.profile_picture, u.is_deleted
		FROM Wallet w
		JOIN users u ON w.user_id = u.user_id
		WHERE w.wallet_id = $1
	`
	row := r.db.QueryRow(query, walletID)
	wallet := &domain.Wallet{}
	user := &domain.User{}
	err := row.Scan(&wallet.ID, &wallet.Balance, &user.ID, &user.Name, &user.Email, &user.Password, &user.ProfilePicture, &user.IsDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Dompet tidak ditemukan
		}
		return nil, fmt.Errorf("failed to find wallet: %v", err)
	}

	wallet.UserId = *user
// Mendapatkan total jumlah top-up yang telah dilakukan untuk wallet ini
// queryTopUp := "SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) FROM TopUp WHERE wallet_id = $1"
// row = r.db.QueryRow(queryTopUp, walletID)
// var totalTopUp float64
// err = row.Scan(&totalTopUp)
// if err != nil {
//     return nil, fmt.Errorf("failed to get total top-up amount: %v", err)
// }
// // Menambahkan total top-up ke balance wallet
// wallet.Balance=0
// wallet.Balance += totalTopUp
	return wallet, nil
}




func (r *walletRepository) Update(wallet *domain.Wallet) error {
	// Mengambil saldo top-up terakhir untuk wallet ini
	queryTopUp := "SELECT amount FROM TopUp WHERE wallet_id = $1 ORDER BY timestamp DESC LIMIT 1"
	row := r.db.QueryRow(queryTopUp, wallet.ID)
	var topupAmount float64
	err := row.Scan(&topupAmount)
	if err != nil {
		return fmt.Errorf("failed to get top-up amount: %v", err)
	}
	// Menambahkan saldo top-up ke balance wallet
	wallet.Balance += topupAmount
	// Perbarui nilai balance pada tabel wallet dengan nilai baru
	queryUpdate := "UPDATE Wallet SET balance = $1 WHERE wallet_id = $2"
	_, err = r.db.Exec(queryUpdate, wallet.Balance, wallet.ID)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %v", err)
	}
	return nil
}

func (r *walletRepository) Delete(walletID int) error {
	query := "DELETE FROM Wallet WHERE wallet_id = $1"
	_, err := r.db.Exec(query, walletID)
	if err != nil {
		return fmt.Errorf("failed to delete wallet: %v", err)
	}
	return nil
}
