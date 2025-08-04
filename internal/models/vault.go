package models

import "time"

var (
	VaultTypePassword = "password"
	VaultTypeNote     = "note"
	VaultTypeCard     = "card"
	VaultTypeBinary   = "binary"
)

type VaultCreateItemReq struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	EncryptedData string `json:"encrypted_data"`
}

type VaultItem struct {
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
	DecryptedData       any       `db:"-" json:"-"`
	Name                string    `db:"name" json:"name"`
	Type                string    `db:"type" json:"type"`
	BinaryPath          string    `db:"binary_path" json:"-"`
	EncryptedDataString string    `db:"-" json:"encrypted_data"`
	EncryptedData       []byte    `db:"encrypted_data" json:"-"`
	ID                  int       `db:"id" json:"id"`
	UserID              int       `db:"user_id" json:"user_id"`
}

type VaultItemState struct {
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
	DecryptedPassword VaultItemDataPassword `json:"-"`
	DecryptedCard     VaultItemDataCard     `json:"-"`
	Name              string                `json:"name"`
	Type              string                `json:"type"`
	DecryptedNote     VaultNote             `json:"-"`
	EncryptedData     []byte                `json:"-"`
	ID                int                   `json:"id"`
}

type VaultNote struct {
	Text string `json:"text"`
}

type VaultItemDataPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

type VaultItemDataCard struct {
	Number string `json:"number"`
	Exp    string `json:"exp"`
	CVV    string `json:"cvv"`
}

type VaultBinaryItemUploadReq struct {
	Name          string `json:"name"`
	EncryptedData string `json:"encrypted_data"`
}

type VaultBinaryItemUploadRes struct {
	URL    string    `json:"url"`
	Path   string    `json:"path"`
	Item   VaultItem `json:"item"`
	ItemID int       `json:"vault_id"`
}

type VaultConfirmBinaryUploadReq struct {
	Path    string `json:"path" db:"path"`
	VaultID int    `json:"vault_id" db:"vault_id"`
}

type VaultItemDataBinary struct {
	Path string `json:"name"`
	Size int64  `json:"size"`
}

type GetBinaryFileURLRes struct {
	URL string `json:"url"`
}
