package errs

import "errors"

var (
	ErrNotValidContextLogin = errors.New("not valid context login")
	ErrInvalidJWTToken      = errors.New("invalid token claims or token not valid")
)
