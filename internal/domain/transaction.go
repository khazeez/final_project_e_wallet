package domain

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	WalletId  Wallet    `json:"wallet_id"`
	Type      string    `json:"type"`
	Amount    int64     `json:"amount"`
	Timestamp time.Time `json:"time"`
}
