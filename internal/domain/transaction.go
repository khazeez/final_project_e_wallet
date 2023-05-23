package domain

import "time"

type Transaction struct {
	ID int
	WalletId Wallet
	Type string
	Amount int64
	Timestamp time.Time
}