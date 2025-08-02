package errs

import (
	"errors"
)

// General errors.
var (
	ErrInvalidType       = errors.New("invalid type")
	ErrNotValidContextID = errors.New("not valid context id")
	ErrValueIsNotString  = errors.New("value is not a string")
)

// Client errors.
var (
	ErrMustContainHTTP = errors.New("must contain http:// or https://")
)

// Repository errors.
var (
	ErrRepoRowsNotFound = errors.New("rows not found")
)

// Auth errors.
var (
	ErrLoginAlreadyExists     = errors.New("login already exist")
	ErrInvalidCredentials     = errors.New("invalid login or password")
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrEqualPasswords         = errors.New("password and crypto password must not be equal")
	ErrInvalidCryptoPasssword = errors.New("invalid crypto password")
)

// Crypto errors.
var (
	ErrSaltNotValid     = errors.New("invalid salt")
	ErrSaltMustBase64   = errors.New("salt must be encoded with base64")
	ErrSaltNotValidLen  = errors.New("invalid salt")
	ErrInvalidEncrypted = errors.New("invalid encrypted data")
	ErrEmptySalt        = errors.New("salt must not empty")
)

// Vault errors.
var (
	ErrVaultNotCorrectType = errors.New("type is not valid")
	ErrBinaryFileNotFound  = errors.New("binary file not found")
	ErrNotFourndActiveItem = errors.New("not found active item")
	ErrFileEmpty           = errors.New("file must not empty")
	ErrBinaryFileUpload    = errors.New("binary file upload err")
)
