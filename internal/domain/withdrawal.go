package domain

import "time"

type Withdrawal struct {
	ID int
	WalletId Wallet
	Amount int64
	Timestamp time.Time
}