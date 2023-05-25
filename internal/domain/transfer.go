package domain

import (
	"time"
)

type Transfer struct {
	ID int `json:"id"`
	SenderId Wallet `json:"sender_id"`
	ReceiferId Wallet `json:"receifer_id"`
	Amount float64 `json:"amount"`
	Timestamp time.Time `json:"time"`
}