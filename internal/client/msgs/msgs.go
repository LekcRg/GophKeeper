package msgs

import "github.com/LekcRg/GophKeeper/internal/models"

type ErrorMsg error

type SelectAuthMsg struct {
	Address  string
	Selected string
}

type FormSubmitMsg struct {
	Values     map[string]string
	ButtonName string
}

type CredentialsBytesMsg struct {
	Key  string
	Salt []byte
	Tag  []byte
}

type UpdateKeySuccessMsg string

type CryptoPassValid string

type ListLoaded struct{}

type SelectVaultItem int

type ToCreateVaultItem struct{}

type SelectTypeMsg struct {
	Selected string
}

type CreateVaultSuccess struct {
	Item models.VaultItem
}

type UpdateAndSwitchToTable struct{}
