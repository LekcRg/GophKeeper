package models

type User struct {
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"-" db:"passhash"`
	KeyHash      string `json:"-" db:"key_hash"`
	EncryptedTag string `json:"encrypted_tag" db:"encrypted_tag"`
	Salt         string `json:"salt" db:"salt"`
	ID           int    `json:"id" db:"id"`
}

type UserReq struct {
	Login        string `json:"login" db:"login"`
	Password     string `json:"password" db:"-"`
	PasswordHash string `json:"-" db:"passhash"`
	KeyHash      string `json:"-" db:"key_hash"`
	EncryptedTag string `json:"encrypted_tag" db:"encrypted_tag"`
	Salt         string `json:"salt" db:"salt"`
}

type UserChangePasswordReq struct {
	Login           string `json:"login" db:"login"`
	CurrentPassword string `json:"current-password" db:"-"`
	NewPassword     string `json:"new-password" db:"-"`
}

type APIKeyRes struct {
	Key string `json:"key"`
}

type UserLogin struct {
	Login        string `json:"login" db:"login"`
	Password     string `json:"password" db:"-"`
	PasswordHash string `json:"-"`
}

type CryptoParamsRes struct {
	EncryptedTag string `json:"encrypted_tag"`
	Salt         string `json:"salt"`
}
