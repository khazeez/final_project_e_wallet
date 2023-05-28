package domain

import "time"

type Payment struct {
	ID            int       `json:"id"`
	WalletId      Wallet    `json:"wallet_id"`
	Amount        float64   `json:"amount"`
	Timestamp     time.Time `json:"time"`
	PaymentType   string    `json:"payment_type"`
	PaymentDetail string    `json:"payment_detail"`
}
