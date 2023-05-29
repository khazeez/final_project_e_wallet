package domain

import (
	"time"
)

type Transfer struct {
	ID         int            `json:"id"`
	SenderId   SenderWallet   `json:"sender_id"`
	ReceiferId ReceiverWallet `json:"receifer_id"`
	Amount     float64        `json:"amount"`
	Timestamp  time.Time      `json:"time"`
}
