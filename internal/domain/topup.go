package domain

import "time"

type TopUp struct {
	ID int
	WalletId Wallet
	Amount int64
	Timestamp time.Time
}