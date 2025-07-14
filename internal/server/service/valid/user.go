package valid

import (
	"regexp"

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

func Register(user *models.UserReq) error {
	return validation.ValidateStruct(user,
		loginField(&user.Login),
		passwordField(&user.Password),
	)
}

func ChangePassword(user *models.UserChangePasswordReq) error {
	return validation.ValidateStruct(user,
		passwordField(&user.NewPassword),
	)
}
