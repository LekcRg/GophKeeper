package models

type User struct {
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"-" db:"passhash"`
	ID           int    `json:"id" db:"id"`
}

type UserReq struct {
	Login        string `json:"login" db:"login"`
	Password     string `json:"password" db:"-"`
	PasswordHash string `json:"-" db:"passhash"`
}

type TokenUserRes struct {
	Token string `json:"token"`
}

type UserChangePasswordReq struct {
	Login           string `json:"-"`
	CurrentPassword string `json:"current-password" db:"-"`
	NewPassword     string `json:"new-password" db:"-"`
}
