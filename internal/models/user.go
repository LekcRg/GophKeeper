package models

type CreateUserReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
}
