package domain

type Wallet struct {
	ID      int     `json:"id"`
	UserId  User    `json:"user_id"`
	Balance float64 `json:"balance"`
}

type SenderWallet struct {
	ID      int        `json:"id"`
	UserId  UserSender `json:"user_sender_id"`
	Balance float64    `json:"balance"`
}

type ReceiverWallet struct {
	ID      int         `json:"id"`
	UserId  UserReceiver `json:"user_receiver_id"`
	Balance float64      `json:"balance"`
}
