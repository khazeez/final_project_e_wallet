package repository

import (
	"database/sql"
	"fmt"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)

type TransactionRepository interface {
	CreateTransaction(transaction *domain.Transaction) error
	GetTransactionsByWalletID(walletID int) ([]*domain.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) CreateTransaction(transaction *domain.Transaction) error {
	query := "INSERT INTO Transaction (transaction_id, wallet_id, type, amount, timestamp) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, transaction.ID, transaction.WalletId, transaction.Type, transaction.Amount, transaction.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}
	return nil
}

func (r *transactionRepository) GetTransactionsByWalletID(walletID int) ([]*domain.Transaction, error) {
	query := "SELECT transaction_id, wallet_id, type, amount, timestamp FROM Transaction WHERE wallet_id = $1"
	rows, err := r.db.Query(query, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}
	defer rows.Close()

	transactions := make([]*domain.Transaction, 0)
	for rows.Next() {
		transaction := &domain.Transaction{}
		err := rows.Scan(&transaction.ID, &transaction.WalletId, &transaction.Type, &transaction.Amount, &transaction.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("unexpected error occurred while iterating rows: %v", err)
	}

	return transactions, nil
}
