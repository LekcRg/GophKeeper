package msgs

import "github.com/LekcRg/GophKeeper/internal/models"

type RegisterSuccessMsg struct {
	Res models.ClientRegisterResponse
}

type RegisterErrorMsg struct {
	Err error
}
