package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)

type TopupRepository interface {
	Create(wallet *domain.TopUp) error
	FindOne(walletID int) (*domain.TopUp, error)
	Update(*domain.TopUp) error
	Delete(walletID int) error
	HistoryTopup(wallet_id int) ([]*domain.TopUp, error)
	GetLastTopupAmount(walletID int) (float64, error)
}

type topupRepository struct {
	db *sql.DB
}

func NewTopupRepository(db *sql.DB) TopupRepository {
	return &topupRepository{
		db: db,
	}
}

func (r *topupRepository) Create(topup *domain.TopUp) error {
	// Update saldo pada tabel Wallet
	updateQuery := "UPDATE Wallet SET balance = balance + $1 WHERE wallet_id = $2"
	_, err := r.db.Exec(updateQuery, topup.Amount, topup.WalletId.ID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}

	// Simpan data top-up ke dalam tabel TopUp
	time := time.Now()
	insertQuery := "INSERT INTO TopUp (topup_id, wallet_id, amount, timestamp) VALUES ($1, $2, $3, $4)"
	_, err = r.db.Exec(insertQuery, topup.ID, topup.WalletId.ID, topup.Amount, time)
	if err != nil {
		// Jika terjadi kesalahan saat menyimpan top-up, perlu mengembalikan saldo ke kondisi semula
		revertQuery := "UPDATE Wallet SET balance = balance - $1 WHERE wallet_id = $2"
		_, revertErr := r.db.Exec(revertQuery, topup.Amount, topup.WalletId.ID)
		if revertErr != nil {
			return fmt.Errorf("failed to create top-up and revert wallet balance: %v (revert error: %v)", err, revertErr)
		}
		return fmt.Errorf("failed to create top-up: %v", err)
	}

	return nil
}

func (r *topupRepository) GetLastTopupAmount(walletID int) (float64, error) {
	query := "SELECT amount FROM TopUp WHERE wallet_id = $1 ORDER BY topup_id DESC LIMIT 1"
	row := r.db.QueryRow(query, walletID)
	var topupAmount float64
	err := row.Scan(&topupAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no top-up found for the wallet")
		}
		return 0, fmt.Errorf("failed to get last top-up amount: %v", err)
	}
	return topupAmount, nil
}

func (r *topupRepository) FindOne(topupID int) (*domain.TopUp, error) {
	query := `
		SELECT t.topup_id, t.amount, t.timestamp, w.wallet_id, u.user_id, u.name, u.email, u.password, u.profile_picture, u.is_deleted, w.balance
		FROM TopUp t
		JOIN Wallet w ON t.wallet_id = w.wallet_id
		JOIN users u ON w.user_id = u.user_id
		WHERE t.topup_id = $1
	`

	row := r.db.QueryRow(query, topupID)
	topup := &domain.TopUp{}
	wallet := &domain.Wallet{}
	user := &domain.User{}

	err := row.Scan(
		&topup.ID,
		&topup.Amount,
		&topup.Timestamp,
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
			return nil, fmt.Errorf("top-up not found")
		}
		return nil, fmt.Errorf("failed to find top-up: %v", err)
	}

	// Set user to wallet's UserId
	wallet.UserId = *user
	// Set wallet to topup's WalletId
	topup.WalletId = *wallet
	return topup, nil
}

func (r *topupRepository) Update(topup *domain.TopUp) error {
	// Mendapatkan nilai jumlah top-up sebelumnya
	queryGetAmount := "SELECT amount FROM topup WHERE topup_id = $1"
	row := r.db.QueryRow(queryGetAmount, topup.ID)
	var previousAmount float64
	err := row.Scan(&previousAmount)
	if err != nil {
		return fmt.Errorf("failed to get previous top-up amount: %v", err)
	}

	// Menghitung jumlah baru dengan menambahkan nilai sebelumnya
	newAmount := 0.00
	newAmount = previousAmount + topup.Amount

	// Memperbarui jumlah top-up dengan nilai baru
	queryUpdateAmount := "UPDATE topup SET amount = $1 WHERE topup_id = $2"
	_, err = r.db.Exec(queryUpdateAmount, newAmount, topup.ID)
	if err != nil {
		// Jika terjadi kesalahan saat menyimpan top-up, perlu mengembalikan saldo ke kondisi semula
		revertQuery := "UPDATE Wallet SET balance = balance - $1 WHERE wallet_id = $2"
		_, revertErr := r.db.Exec(revertQuery, topup.Amount, topup.WalletId.ID)
		if revertErr != nil {
			return fmt.Errorf("failed to create top-up and revert wallet balance: %v (revert error: %v)", err, revertErr)
		}
		return fmt.Errorf("failed to create top-up: %v", err)
	}

	return nil
}

func (r *topupRepository) Delete(topupID int) error {
	query := "DELETE FROM TopUp WHERE topup_id = $1"
	_, err := r.db.Exec(query, topupID)
	if err != nil {
		return fmt.Errorf("failed to delete top-up: %v", err)
	}
	return nil
}

func (r *topupRepository) HistoryTopup(wallet_id int) ([]*domain.TopUp, error) {
	query := `SELECT t.topup_id, t.amount,t.timestamp, w.wallet_id, u.user_id, u.name, u.email, u.password, u.profile_picture, u.is_deleted, w.balance FROM topup t
	 	JOIN Wallet w ON t.wallet_id = w.wallet_id
	 	JOIN users u ON w.user_id = u.user_id
	 	WHERE t.wallet_id = $1`

	rows, err := r.db.Query(query, wallet_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawals: %v", err)
	}
	defer rows.Close()

	topups := []*domain.TopUp{}
	wallet := &domain.Wallet{}
	user := &domain.User{}

	for rows.Next() {
		topup := &domain.TopUp{}
		err := rows.Scan(
	&topup.ID,
	&topup.Amount,
	&topup.Timestamp,
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
			return nil, fmt.Errorf("failed to scan withdrawal row: %v", err)
		}

	wallet.UserId = *user
	topup.WalletId = *wallet

	topups = append(topups, topup)
	}
	return topups, nil
}


