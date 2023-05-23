package domain

type Wallet struct {
	ID int
	UserId User
	Balance int64
}