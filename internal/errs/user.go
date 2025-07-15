package errs

import "errors"

var (
	ErrLoginAlreadyExists    = errors.New("login already exist")
	ErrInvalidCredentials    = errors.New("invalid login or password")
	ErrUserWithLoginNotFound = errors.New("user with this login not found")
	ErrInvalidPassword       = errors.New("invalid password")
)
