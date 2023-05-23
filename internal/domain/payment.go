package domain

import "time"

type Payment struct {
	ID int
	WalletId int
	Amount int64
	Timestamp time.Time
	PaymentType string
	PaymentDetail string
}