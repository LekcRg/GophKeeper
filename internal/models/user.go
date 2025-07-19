package models

type User struct {
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"-" db:"passhash"`
	KeyHash      string `json:"-" db:"key_hash"`
	ID           int    `json:"id" db:"id"`
}

type UserReq struct {
	Login        string `json:"login" db:"login"`
	Password     string `json:"password" db:"-"`
	PasswordHash string `json:"-" db:"passhash"`
	KeyHash      string `json:"-" db:"key_hash"`
}

type UserChangePasswordReq struct {
	Login           string `json:"login" db:"login"`
	CurrentPassword string `json:"current-password" db:"-"`
	NewPassword     string `json:"new-password" db:"-"`
}

type APIKeyRes struct {
	Key string `json:"key"`
}
