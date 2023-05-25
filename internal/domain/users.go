package domain

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profile_picture"`
	IsDeleted      bool   `json:"is_deleted"`
}


type LoginUser struct {
	Username string	`json:"username"`
	Password string `json:"password"`
}
