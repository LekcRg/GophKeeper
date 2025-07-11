package errs

import "errors"

var (
	ErrLoginAlreadyExists = errors.New("login already exist")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrInvalidPassword    = errors.New("invalid password")
)
