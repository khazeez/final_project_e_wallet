package domain

import (
	"time"
)

type Transfer struct {
	ID int
	SenderId Wallet
	ReceiferId Wallet
	Amount int64
	Timestamp time.Time
}