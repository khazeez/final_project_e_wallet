package domain

import "time"

type Withdrawal struct {
	ID        int       `json:"id"`
	WalletId  Wallet    `json:"wallet_id"`
	Amount    float64     `json:"amount"`
	Timestamp time.Time `json:"time"`
}
