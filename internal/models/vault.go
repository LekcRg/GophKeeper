package models

import "time"

type VaultCreateItemReq struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	EncryptedData string `json:"encrypted_data"`
}

type VaultItem struct {
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	Name          string    `db:"name" json:"name"`
	Type          string    `db:"type" json:"type"`
	EncryptedData []byte    `db:"encrypted_data" json:"-"`
	ID            int       `db:"id" json:"id"`
	UserID        int       `db:"user_id" json:"user_id"`
}
