package valid

import (
	"encoding/base64"
	"errors"
	"maps"
	"slices"
	"sort"
	"testing"

	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	type test struct {
		user       models.UserReq
		name       string
		wantErrs   []string
		checkLogin bool
	}

	saltBytes, err := crypto.GenEncryptionSalt()
	require.NoError(t, err)

	salt := base64.StdEncoding.EncodeToString(saltBytes)

	tagBytes := []byte("tag")
	tag := base64.StdEncoding.EncodeToString(tagBytes)

	tests := []test{
		{
			name:       "success",
			checkLogin: true,
			user: models.UserReq{
				Login:        "user",
				Password:     "P@ssw0rdTesting3",
				Salt:         salt,
				EncryptedTag: tag,
			},
		},
		{
			name: "password without special symbol",
			user: models.UserReq{
				Login:        "user",
				Password:     "Passw0rdTesting3",
				Salt:         salt,
				EncryptedTag: tag,
			},
			wantErrs: []string{"password"},
		},
		{
			name: "password without number",
			user: models.UserReq{
				Login:        "user",
				Password:     "P@sswordTestingq",
				Salt:         salt,
				EncryptedTag: tag,
			},
			wantErrs: []string{"password"},
		},
		{
			name: "password without uppercase letter",
			user: models.UserReq{
				Login:        "user",
				Password:     "p@ssw0rdtesting3",
				Salt:         salt,
				EncryptedTag: tag,
			},
			wantErrs: []string{"password"},
		},
		{
			name: "password with lower len",
			user: models.UserReq{
				Login:        "user",
				Password:     "P@sswrdTt3",
				Salt:         salt,
				EncryptedTag: tag,
			},
			wantErrs: []string{"password"},
		},
		{
			name:       "without password",
			checkLogin: true,
			user: models.UserReq{
				Login:        "user",
				Salt:         salt,
				EncryptedTag: tag,
			},
			wantErrs: []string{"password"},
		},
		{
			name: "login with invalid symbol",
			user: models.UserReq{
				Login:        "user@",
				Password:     "P@ssw0rdTesting3",
				Salt:         salt,
				EncryptedTag: tag,
			},
			wantErrs: []string{"login"},
		},
		{
			name: "login with lower len",
			user: models.UserReq{
				Login:        "usr",
				Password:     "P@ssw0rdTesting3",
				Salt:         salt,
				EncryptedTag: tag,
			},
			wantErrs: []string{"login"},
		},
		{
			name:       "without login",
			checkLogin: true,
			user: models.UserReq{
				Password:     "P@ssw0rdTesting3",
				Salt:         salt,
				EncryptedTag: tag,
			},
			wantErrs: []string{"login"},
		},
		{
			name:       "without tag",
			checkLogin: true,
			user: models.UserReq{
				Login:    "user",
				Password: "P@ssw0rdTesting3",
				Salt:     salt,
			},
			wantErrs: []string{"encrypted_tag"},
		},
		{
			name:       "invalid salt base64",
			checkLogin: true,
			user: models.UserReq{
				Login:        "user",
				Password:     "P@ssw0rdTesting3",
				Salt:         salt[1:],
				EncryptedTag: tag,
			},
			wantErrs: []string{"salt"},
		},
		{
			name:       "invalid salt len",
			checkLogin: true,
			user: models.UserReq{
				Login:        "user",
				Password:     "P@ssw0rdTesting3",
				Salt:         base64.StdEncoding.EncodeToString(saltBytes[1:]),
				EncryptedTag: tag,
			},
			wantErrs: []string{"salt"},
		},
	}

	for _, tt := range tests {
		t.Run("Register_"+tt.name, func(t *testing.T) {
			t.Parallel()

			checkErrorsRun(t, tt.wantErrs, func() error {
				return Register(&tt.user)
			})
		})
	}
}

func TestChangePassword(t *testing.T) {
	t.Parallel()

	type test struct {
		req      models.UserChangePasswordReq
		name     string
		wantErrs []string
	}

	tests := []test{
		{
			name: "success",
			req: models.UserChangePasswordReq{
				Login:           "testuser",
				CurrentPassword: "P@ssw0rdTesting3",
				NewPassword:     "Sup3rSecretP@a$$word",
			},
		},
		{
			name:     "Without login, current-password, new-password",
			req:      models.UserChangePasswordReq{},
			wantErrs: []string{"current-password", "login", "new-password"},
		},
		{
			name: "Without current-password, new-password",
			req: models.UserChangePasswordReq{
				Login: "testuser",
			},
			wantErrs: []string{"current-password", "new-password"},
		},
		{
			name: "Without new-password",
			req: models.UserChangePasswordReq{
				Login:           "testuser",
				CurrentPassword: "Sup3rSecretP@a$$word",
			},
			wantErrs: []string{"new-password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			checkErrorsRun(t, tt.wantErrs, func() error {
				return ChangePassword(&tt.req)
			})
		})
	}
}

func checkErrorsRun(t *testing.T, wantErrs []string, f func() error) {
	t.Helper()

	err := f()

	if len(wantErrs) == 0 {
		require.NoError(t, err)

		return
	}

	var validErr validation.Errors

	ok := errors.As(err, &validErr)
	require.True(t, ok)

	keyErrs := slices.Collect(maps.Keys(validErr))
	sort.Strings(keyErrs)

	t.Log(keyErrs)
	assert.Equal(t, wantErrs, keyErrs)
}

func TestLogin(t *testing.T) {
	t.Parallel()

	type test struct {
		user       models.UserLogin
		name       string
		wantErrs   []string
		checkLogin bool
	}

	tests := []test{
		{
			name:       "success",
			checkLogin: true,
			user: models.UserLogin{
				Login:    "user",
				Password: "P@ssw0rdTesting3",
			},
		},
		{
			name:       "without password",
			checkLogin: true,
			user: models.UserLogin{
				Login: "user",
			},
			wantErrs: []string{"password"},
		},
		{
			name:       "without login",
			checkLogin: true,
			user: models.UserLogin{
				Password: "P@ssw0rdTesting3",
			},
			wantErrs: []string{"login"},
		},
	}

	for _, tt := range tests {
		t.Run("Register_"+tt.name, func(t *testing.T) {
			t.Parallel()

			checkErrorsRun(t, tt.wantErrs, func() error {
				return Login(&tt.user)
			})
		})
	}
}
