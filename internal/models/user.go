package models

type RegisterUserReq struct {
	Login        string `json:"login" db:"login"`
	Password     string `json:"password" db:"-"`
	PasswordHash string `json:"-" db:"passhash"`
}

type RegisterUserRes struct {
	Token string `json:"token"`
}

type RegisterUserError struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type User struct {
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"-" db:"passhash"`
	ID           int    `json:"id" db:"id"`
}
