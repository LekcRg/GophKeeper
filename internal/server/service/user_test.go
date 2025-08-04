package service

import (
	"context"
	"encoding/base64"
	"errors"
	"math/rand"
	"testing"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/mocks"
	"github.com/LekcRg/GophKeeper/internal/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	cfg, err := config.GetConfig([]string{})
	require.NoError(t, err)

	type test struct {
		mockErr   error
		req       models.UserReq
		name      string
		doNotMock bool
		wantErr   bool
	}

	saltBytes, err := crypto.GenEncryptionSalt()
	require.NoError(t, err)

	salt := base64.StdEncoding.EncodeToString(saltBytes)

	tests := []test{
		{
			name: "success",
			req: models.UserReq{
				Login:        "testuser",
				Password:     "T3stP@as5word",
				EncryptedTag: "1234",
				Salt:         salt,
			},
		},
		{
			name: "without password",
			req: models.UserReq{
				Login:        "testuser",
				EncryptedTag: "1234",
				Salt:         salt,
			},
			wantErr:   true,
			doNotMock: true,
		},
		{
			name: "login already exist",
			req: models.UserReq{
				Login:        "testuser",
				Password:     "T3stP@as5word",
				EncryptedTag: "1234",
				Salt:         salt,
			},
			mockErr: errs.ErrLoginAlreadyExists,
			wantErr: true,
		},
		{
			name: "unexpected repository error",
			req: models.UserReq{
				Login:        "testuser",
				Password:     "T3stP@as5word",
				EncryptedTag: "1234",
				Salt:         salt,
			},
			mockErr: &pgconn.PgError{Code: pgerrcode.CannotConnectNow},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockUserRepo(t)
			us := NewUserService(repo, cfg)

			userID := rand.Int()
			if !tt.doNotMock {
				repo.EXPECT().
					CreateUser(mock.Anything, mock.Anything).
					Return(userID, tt.mockErr)
			}

			res, err := us.Register(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.NotEmpty(t, res.Key)
		})
	}
}

func TestUpdateAPIKey(t *testing.T) {
	t.Parallel()

	cfg, err := config.GetConfig([]string{})
	require.NoError(t, err)

	password := "password"
	passwordHash, err := crypto.HashPassword(password)
	require.NoError(t, err)

	type test struct {
		mockErr         error
		req             models.UserLogin
		name            string
		doNotMock       bool
		doNotMockUpdate bool
		GetUserErr      error
		UpdateKeyErr    error
		isErr           error
		asErr           any
	}

	errInternal := errors.New("internal")
	tests := []test{
		{
			name: "success",
			req: models.UserLogin{
				Login:    "login",
				Password: password,
			},
		},
		{
			name: "Empty password",
			req: models.UserLogin{
				Login: "login",
			},
			asErr:     validation.Errors{},
			doNotMock: true,
		},
		{
			name: "Invalid password",
			req: models.UserLogin{
				Login:    "login",
				Password: "invalid",
			},
			isErr:           errs.ErrInvalidCredentials,
			doNotMockUpdate: true,
		},
		{
			name: "Not found user",
			req: models.UserLogin{
				Login:    "login",
				Password: "password",
			},
			isErr:           errs.ErrInvalidCredentials,
			GetUserErr:      errs.ErrUserNotFound,
			doNotMockUpdate: true,
		},
		{
			name: "Internal get user error",
			req: models.UserLogin{
				Login:    "login",
				Password: "password",
			},
			isErr:           errInternal,
			GetUserErr:      errInternal,
			doNotMockUpdate: true,
		},
		{
			name: "UpdateUserKey internal error",
			req: models.UserLogin{
				Login:    "login",
				Password: "password",
			},
			isErr:        errInternal,
			UpdateKeyErr: errInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockUserRepo(t)
			us := NewUserService(repo, cfg)

			if !tt.doNotMock {
				repo.EXPECT().
					GetUserByLogin(mock.Anything, tt.req.Login).
					Return(models.User{
						Login:        tt.req.Login,
						PasswordHash: passwordHash,
					}, tt.GetUserErr)

				if !tt.doNotMockUpdate {
					repo.EXPECT().
						UpdateUserKey(mock.Anything, mock.Anything).
						Return(tt.UpdateKeyErr)
				}
			}

			res, err := us.UpdateAPIKey(context.Background(), tt.req)
			if tt.isErr != nil || tt.asErr != nil {
				assert.Error(t, err)

				if tt.isErr != nil {
					assert.ErrorIs(t, err, tt.isErr)
				}

				if tt.asErr != nil {
					assert.ErrorAs(t, err, &tt.asErr)
				}

				return
			}

			assert.NotNil(t, res.Key)
		})
	}
}

func TestChangePassword(t *testing.T) {
	t.Parallel()

	cfg, err := config.GetConfig([]string{})
	require.NoError(t, err)

	password := "T3stP@as5word"
	newPassword := "N3wSup3rP@55wrd"

	hash, err := crypto.HashPassword(password)
	require.NoError(t, err)

	type test struct {
		mockGetErr      error
		mockUpdateErr   error
		checkErr        func(error) bool
		req             models.UserChangePasswordReq
		name            string
		doNotMockGet    bool
		doNotMockUpdate bool
	}

	tests := []test{
		{
			name: "success",
			req: models.UserChangePasswordReq{
				Login:           "testuser",
				CurrentPassword: password,
				NewPassword:     newPassword,
			},
		},
		{
			name: "validation error",
			req: models.UserChangePasswordReq{
				Login:           "testuser",
				CurrentPassword: password,
				NewPassword:     "short",
			},
			checkErr: func(err error) bool {
				var v validation.Errors
				return errors.As(err, &v)
			},
			doNotMockGet: true,
		},
		{
			name: "invalid current password",
			req: models.UserChangePasswordReq{
				Login:           "testuser",
				CurrentPassword: "wrongPassword",
				NewPassword:     newPassword,
			},
			checkErr: func(err error) bool {
				return errors.Is(err, errs.ErrInvalidPassword)
			},
			doNotMockUpdate: true,
		},
		{
			name: "user not found",
			req: models.UserChangePasswordReq{
				Login:           "testuser",
				CurrentPassword: password,
				NewPassword:     newPassword,
			},
			mockGetErr: errs.ErrUserNotFound,
			checkErr: func(err error) bool {
				return errors.Is(err, errs.ErrUserNotFound)
			},
			doNotMockUpdate: true,
		},
		{
			name: "unexpected update error",
			req: models.UserChangePasswordReq{
				Login:           "testuser",
				CurrentPassword: password,
				NewPassword:     newPassword,
			},
			mockUpdateErr: &pgconn.PgError{Code: pgerrcode.CannotConnectNow},
			checkErr: func(err error) bool {
				var pgErr *pgconn.PgError
				return errors.As(err, &pgErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockUserRepo(t)
			us := NewUserService(repo, cfg)

			if !tt.doNotMockGet {
				repo.EXPECT().
					GetUserByLogin(context.Background(), tt.req.Login).
					Return(models.User{
						ID:           1,
						Login:        tt.req.Login,
						PasswordHash: hash,
					}, tt.mockGetErr)
			}

			if !tt.doNotMockGet && !tt.doNotMockUpdate {
				repo.EXPECT().
					UpdateUserPassword(context.Background(), mock.Anything).
					Return(tt.mockUpdateErr)
			}

			err := us.ChangePassword(context.Background(), tt.req)
			if tt.checkErr == nil {
				require.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err))
			}
		})
	}
}

func TestCryptoParams(t *testing.T) {
	type test struct {
		name    string
		id      int
		wantErr bool
		svcErr  error
	}

	tests := []test{
		{
			name: "Success",
			id:   1,
		},
		{
			name:    "Success",
			id:      1,
			wantErr: true,
			svcErr:  errs.ErrUserNotFound,
		},
	}

	cfg, err := config.GetConfig([]string{})
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockUserRepo(t)
			us := NewUserService(repo, cfg)

			repo.EXPECT().
				GetUserByID(mock.Anything, tt.id).
				Return(models.User{
					Salt:         "salt",
					EncryptedTag: "tag",
				}, tt.svcErr)

			params, err := us.GetCryptoParams(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.NotEmpty(t, params.EncryptedTag)
			assert.NotEmpty(t, params.Salt)
		})
	}
}
