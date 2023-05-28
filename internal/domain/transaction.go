package domain

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	WalletId  int    `json:"wallet_id"`
	Type      string    `json:"type"`
	Amount    float64    `json:"amount"`
	Timestamp time.Time `json:"time"`
}
