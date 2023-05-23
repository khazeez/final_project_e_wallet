package domain

type Wallet struct {
	ID      int   `json:"id"`
	UserId  User  `json:"user_id"`
	Balance int64 `json:"balance"`
}
