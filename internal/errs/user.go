package errs

import (
	"errors"
	"fmt"

	"github.com/LekcRg/GophKeeper/internal/crypto"
)

var (
	ErrLoginAlreadyExists = errors.New("login already exist")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrSaltNotValid       = errors.New("invalid salt")
	ErrSaltMustBase64     = errors.New("salt must be ecnoded with base64")
	ErrSaltNotValidLen    = fmt.Errorf("salt must be %d bytes", crypto.SaltLen)
)
