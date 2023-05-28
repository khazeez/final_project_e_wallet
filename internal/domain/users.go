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
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSender struct {
	Sender_ID             int    `json:"sender_id"`
	Sender_Name           string `json:"sender_name"`
	Sender_Email          string `json:"sender_email"`
	Sender_Password       string `json:"sender_password"`
	Sender_ProfilePicture string `json:"profile_picture"`
	IsDeleted             bool   `json:"is_deleted"`
}

type UserReceiver struct {
	Receifer_ID             int    `json:"Receifer_id"`
	Receifer_Name           string `json:"Receifer_name"`
	Receifer_Email          string `json:"Receifer_email"`
	Receifer_Password       string `json:"Receifer_password"`
	Receifer_ProfilePicture string `json:"profile_picture"`
	IsDeleted               bool   `json:"is_deleted"`
}
