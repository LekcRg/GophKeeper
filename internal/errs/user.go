package errs

import "errors"

var (
	ErrLoginAlreadyExists = errors.New("login already exist")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
)
