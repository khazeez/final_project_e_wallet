package repository

import (
	"database/sql"
	"fmt"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)

type WithdrawRepository interface {
	Create(withdrawal *domain.Withdrawal) error
	FindOne(withdrawalID int) (*domain.Withdrawal, error)
	FindAll() ([]*domain.Withdrawal, error)
	Update(withdrawal *domain.Withdrawal) error
	Delete(withdrawalID int) error
	History(wallet_id int) ([]*domain.Withdrawal, error)
}
type withdrawRepository struct {
	db *sql.DB
}

func NewWithdrawRepository(db *sql.DB) WithdrawRepository {
	return &withdrawRepository{
		db: db,
	}
}
func (r *withdrawRepository) Create(withdrawal *domain.Withdrawal) error {
	// Cek apakah wallet dengan ID yang diberikan ada dalam database
	query := "SELECT balance FROM Wallet WHERE wallet_id = $1"
	row := r.db.QueryRow(query, withdrawal.WalletId.ID)
	var balance float64
	err := row.Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("wallet not found")
		}
		return fmt.Errorf("failed to get wallet balance: %v", err)
	}
	// Cek apakah saldo cukup untuk melakukan penarikan
	if balance < float64(withdrawal.Amount) {
		return fmt.Errorf("insufficient balance")
	}

	// Kurangi saldo wallet sesuai dengan jumlah penarikan
	// newBalance := balance - float64(withdrawal.Amount)
	// Update saldo pada tabel wallet
	balance = 0
	updateQuery := "UPDATE Wallet SET balance = balance - $1 WHERE wallet_id = $2"
	_, err = r.db.Exec(updateQuery, withdrawal.Amount, withdrawal.WalletId.ID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}
	// Simpan data withdrawal ke dalam tabel Withdrawal
	insertQuery := "INSERT INTO Withdrawal (withdrawal_id, wallet_id, amount, timestamp) VALUES ($1, $2, $3, $4)"
	_, err = r.db.Exec(insertQuery, withdrawal.ID, withdrawal.WalletId.ID, withdrawal.Amount, withdrawal.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to create withdrawal: %v", err)
	}
	return nil
}

func (r *withdrawRepository) FindOne(withdrawalID int) (*domain.Withdrawal, error) {
	query := `
	SELECT t.withdrawal_id, t.amount, w.wallet_id, u.user_id, u.name, u.email, u.password, u.profile_picture, u.is_deleted, w.balance
	FROM withdrawal t
	JOIN Wallet w ON t.wallet_id = w.wallet_id
	JOIN users u ON w.user_id = u.user_id
	WHERE t.withdrawal_id = $1
	`

	row := r.db.QueryRow(query, withdrawalID)
	withdrawal := &domain.Withdrawal{}
	wallet := &domain.Wallet{}
	user := &domain.User{}

	err := row.Scan(
		&withdrawal.ID,
		&withdrawal.Amount,
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
			return nil, fmt.Errorf("withdrawal not found")
		}
		return nil, fmt.Errorf("failed to find withdrawal: %v", err)
	}

	// Set user to wallet's UserId
	wallet.UserId = *user
	// Set wallet to topup's WalletId
	withdrawal.WalletId = *wallet
	return withdrawal, nil

}

func (r *withdrawRepository) FindAll() ([]*domain.Withdrawal, error) {
	query := "SELECT withdrawal_id, wallet_id, amount, timestamp FROM Withdrawal"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawals: %v", err)
	}
	defer rows.Close()

	withdrawals := []*domain.Withdrawal{}
	for rows.Next() {
		withdrawal := &domain.Withdrawal{}
		err := rows.Scan(&withdrawal.ID, &withdrawal.WalletId, &withdrawal.Amount, &withdrawal.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan withdrawal row: %v", err)
		}
		withdrawals = append(withdrawals, withdrawal)
	}
	return withdrawals, nil
}

func (r *withdrawRepository) Update(withdrawal *domain.Withdrawal) error {
	updateQuery := "UPDATE Withdrawal SET wallet_id = $1 WHERE withdrawal_id = $2"
	_, err := r.db.Exec(updateQuery, withdrawal.WalletId.ID, withdrawal.ID)
	if err != nil {
		return fmt.Errorf("failed to update withdrawal: %v", err)
	}
	return nil
}

func (r *withdrawRepository) Delete(withdrawalID int) error {
	deleteQuery := "DELETE FROM Withdrawal WHERE withdrawal_id = $1"
	_, err := r.db.Exec(deleteQuery, withdrawalID)
	if err != nil {
		return fmt.Errorf("failed to delete withdrawal: %v", err)
	}
	return nil
}

func (r *withdrawRepository) History(wallet_id int) ([]*domain.Withdrawal, error) {
	query := `SELECT t.withdrawal_id, t.amount, w.wallet_id, u.user_id, u.name, u.email, u.password, u.profile_picture, u.is_deleted, w.balance FROM withdrawal t
	 	JOIN Wallet w ON t.wallet_id = w.wallet_id
	 	JOIN users u ON w.user_id = u.user_id
	 	WHERE t.wallet_id = $1`

	rows, err := r.db.Query(query, wallet_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawals: %v", err)
	}
	defer rows.Close()

	withdrawals := []*domain.Withdrawal{}
	wallet := &domain.Wallet{}
	user := &domain.User{}

	for rows.Next() {
		withdrawal := &domain.Withdrawal{}
		err := rows.Scan(
	&withdrawal.ID,
	&withdrawal.Amount,
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
	withdrawal.WalletId = *wallet

	withdrawals = append(withdrawals, withdrawal)
	}
	return withdrawals, nil
}


	