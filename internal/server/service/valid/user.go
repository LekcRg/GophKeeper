package valid

import (
	"encoding/base64"
	"regexp"

	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func passwordField(password *string) *validation.FieldRules {
	const (
		minLenPassword = 12
		maxLenPassword = 50
	)

	return validation.Field(
		password,
		validation.Required,
		validation.Length(minLenPassword, maxLenPassword),
		validation.Match(regexp.MustCompile("[a-z]")).
			Error("password requires lowercase character"),
		validation.Match(regexp.MustCompile("[A-Z]")).
			Error("password requires uppercase character"),
		validation.Match(regexp.MustCompile("[1-9]")).
			Error("password requires number"),
		validation.Match(regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`)).
			Error("password requires special character"),
	)
}

func loginField(login *string) *validation.FieldRules {
	const (
		minLenLogin = 4
		maxLenLogin = 50
	)

	return validation.Field(
		login,
		validation.Required,
		validation.Length(minLenLogin, maxLenLogin),
		validation.Match(regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)).
			Error("login must contain only letters, numbers, underscores, or hyphens"),
	)
}

func salt(value any) error {
	b64salt, ok := value.(string)
	if !ok {
		return errs.ErrSaltNotValid
	}

	salt, err := base64.StdEncoding.DecodeString(b64salt)
	if err != nil {
		return errs.ErrSaltMustBase64
	}

	if len(salt) != crypto.SaltLen {
		return errs.ErrSaltNotValidLen
	}

	return nil
}

func Register(user *models.UserReq) error {
	return validation.ValidateStruct(user,
		loginField(&user.Login),
		passwordField(&user.Password),
		validation.Field(&user.Salt,
			validation.Required,
			validation.By(salt)),
		validation.Field(&user.EncryptedTag, validation.Required),
	)
}

func Login(user *models.UserLogin) error {
	return validation.ValidateStruct(user,
		validation.Field(&user.Login, validation.Required),
		validation.Field(&user.Password, validation.Required),
	)
}

func ChangePassword(user *models.UserChangePasswordReq) error {
	return validation.ValidateStruct(user,
		validation.Field(&user.Login, validation.Required),
		validation.Field(&user.CurrentPassword, validation.Required),
		passwordField(&user.NewPassword),
	)
}
