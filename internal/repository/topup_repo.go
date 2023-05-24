package repository

import (
	"database/sql"
	"fmt"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)

type TopupRepository interface {
	Create(wallet *domain.TopUp) error
	FindOne(walletID int) (*domain.TopUp, error)
	Update(*domain.TopUp) error 
	Delete(walletID int) error
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
	query := "INSERT INTO TopUp (topup_id, amount) VALUES (?, ?)"
	_, err := r.db.Exec(query, topup.ID, topup.Amount)
	if err != nil {
		return fmt.Errorf("failed to create top-up: %v", err)
	}
	return nil
}

func (r *topupRepository) GetLastTopupAmount(walletID int) (float64, error) {
	query := "SELECT amount FROM TopUp WHERE wallet_id = ? ORDER BY topup_id DESC LIMIT 1"
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
	query := "SELECT topup_id, amount FROM TopUp WHERE topup_id = ?"
	row := r.db.QueryRow(query, topupID)
	topup := &domain.TopUp{}
	err := row.Scan(&topup.ID, &topup.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("top-up not found")
		}
		return nil, fmt.Errorf("failed to find top-up: %v", err)
	}
	return topup, nil
}
func (r *topupRepository) Update(topup *domain.TopUp) error {
	query := "UPDATE TopUp SET amount = ? WHERE topup_id = ?"
	_, err := r.db.Exec(query, topup.Amount, topup.ID)
	if err != nil {
		return fmt.Errorf("failed to update top-up: %v", err)
	}
	return nil
}

func (r *topupRepository) Delete(topupID int) error {
	query := "DELETE FROM TopUp WHERE topup_id = ?"
	_, err := r.db.Exec(query, topupID)
	if err != nil {
		return fmt.Errorf("failed to delete top-up: %v", err)
	}
	return nil
}


