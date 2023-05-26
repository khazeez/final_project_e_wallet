package domain

type Wallet struct {
	ID      int    `json:"id"`
	UserId  User    `json:"user_id"`
	Balance float64 `json:"balance"`
}
