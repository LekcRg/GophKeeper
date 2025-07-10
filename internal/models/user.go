package models

type User struct {
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"-" db:"passhash"`
	ID           int    `json:"id" db:"id"`
}

type RegisterUserReq struct {
	Login        string `json:"login" db:"login"`
	Password     string `json:"password" db:"-"`
	PasswordHash string `json:"-" db:"passhash"`
}

type UserError struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginUserReq struct {
	Login        string `json:"login" db:"login"`
	Password     string `json:"password" db:"-"`
	PasswordHash string `json:"-" db:"passhash"`
}

type TokenUserRes struct {
	Token string `json:"token"`
}
