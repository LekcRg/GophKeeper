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

type RegisterSuccessMsg struct {
	Res models.ClientRegisterResponse
}

type UpdateKeySuccessMsg string
