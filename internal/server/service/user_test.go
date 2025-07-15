package service

import (
	"context"
	"errors"
	"testing"
	"time"

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

var cfg = &config.Config{Auth: config.Auth{
	Secret:    "testsecret",
	JWTExpire: time.Minute * 5,
}}

type testUserService func(ctx context.Context, req models.UserReq) (models.TokenUserRes, error)

type testLoginRegister struct {
	mockErr   error
	checkErr  func(err error) bool
	req       models.UserReq
	name      string
	doNotMock bool
}

func TestRegister(t *testing.T) {
	t.Parallel()

	tests := []testLoginRegister{
		{
			name: "success",
			req: models.UserReq{
				Login:    "testuser",
				Password: "T3stP@as5word",
			},
		},
		{
			name: "without password",
			req: models.UserReq{
				Login: "testuser",
			},
			checkErr: func(err error) bool {
				var validErr validation.Errors

				return errors.As(err, &validErr)
			},
			doNotMock: true,
		},
		{
			name: "login already exist",
			req: models.UserReq{
				Login:    "testuser",
				Password: "T3stP@as5word",
			},
			checkErr: func(err error) bool {
				return errors.Is(err, errs.ErrLoginAlreadyExists)
			},
			mockErr: errs.ErrLoginAlreadyExists,
		},
		{
			name: "unexpected repository error",
			req: models.UserReq{
				Login:    "testuser",
				Password: "T3stP@as5word",
			},
			checkErr: func(err error) bool {
				var pgErr *pgconn.PgError

				return errors.As(err, &pgErr)
			},
			mockErr: &pgconn.PgError{
				Code: pgerrcode.CannotConnectNow,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testLoginRegisterRun(t, &tt, func(repo *mocks.MockUserRepo) {
				repo.EXPECT().
					CreateUser(mock.Anything, mock.Anything).
					Return(tt.mockErr)
			}, func(us *UserService) testUserService {
				return us.Register
			})
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	const password = "T3stP@as5word"
	hash, err := crypto.HashPassword(password)
	require.NoError(t, err)

	tests := []testLoginRegister{
		{
			name: "success",
			req: models.UserReq{
				Login:    "testuser",
				Password: password,
			},
		},
		{
			name: "without password",
			req: models.UserReq{
				Login: "testuser",
			},
			checkErr: func(err error) bool {
				var validErr validation.Errors

				return errors.As(err, &validErr)
			},
			doNotMock: true,
		},
		{
			name: "Invalid password",
			req: models.UserReq{
				Login:    "testuser",
				Password: "Test password 123",
			},
			checkErr: func(err error) bool {
				return errors.Is(err, errs.ErrInvalidCredentials)
			},
		},
		{
			name: "Invalid password",
			req: models.UserReq{
				Login:    "testuser",
				Password: password,
			},
			mockErr: errs.ErrUserWithLoginNotFound,
			checkErr: func(err error) bool {
				return errors.Is(err, errs.ErrInvalidCredentials)
			},
		},
		{
			name: "unexpected repository error",
			req: models.UserReq{
				Login:    "testuser",
				Password: password,
			},
			checkErr: func(err error) bool {
				var pgErr *pgconn.PgError

				return errors.As(err, &pgErr)
			},
			mockErr: &pgconn.PgError{
				Code: pgerrcode.CannotConnectNow,
			},
		},
	}

	for _, tt := range tests {
		testLoginRegisterRun(t, &tt, func(repo *mocks.MockUserRepo) {
			repo.EXPECT().
				GetUserByLogin(mock.Anything, tt.req.Login).
				Return(models.User{
					Login:        tt.req.Login,
					PasswordHash: hash,
				}, tt.mockErr)
		}, func(us *UserService) testUserService {
			return us.Login
		})
	}
}

func testLoginRegisterRun(
	t *testing.T, tt *testLoginRegister,
	mockFunc func(repo *mocks.MockUserRepo),
	svcFunc func(us *UserService) testUserService,
) {
	t.Helper()

	repo := mocks.NewMockUserRepo(t)
	us := NewUserService(repo, cfg)

	if !tt.doNotMock {
		mockFunc(repo)
	}

	tokenRes, err := svcFunc(us)(context.Background(), tt.req)
	if tt.checkErr == nil {
		require.NoError(t, err)
		assert.NotEmpty(t, tokenRes.Token)
	} else {
		assert.Error(t, err)
		assert.True(t, tt.checkErr(err))
	}
}

func TestChangePassword(t *testing.T) {
	t.Parallel()

	password := "T3stP@as5word"
	hash, err := crypto.HashPassword(password)
	require.NoError(t, err)

	newPassword := "N3wSup3rP@55wrd"

	type test struct {
		mockGetErr      error
		mockUpdateErr   error
		checkErr        func(err error) bool
		req             models.UserChangePasswordReq
		name            string
		doNotMock       bool
		doNotMockUpdate bool
	}

	tests := []test{
		{
			name: "success",
			req: models.UserChangePasswordReq{
				CurrentPassword: password,
				NewPassword:     newPassword,
				Login:           "testuser",
			},
		},
		{
			name: "validation error",
			req: models.UserChangePasswordReq{
				CurrentPassword: password,
				NewPassword:     "superPassword12345",
				Login:           "testuser",
			},
			checkErr: func(err error) bool {
				var validErr validation.Errors

				return errors.As(err, &validErr)
			},
			doNotMock: true,
		},
		{
			name: "invalid password",
			req: models.UserChangePasswordReq{
				CurrentPassword: "invalid",
				NewPassword:     newPassword,
				Login:           "testuser1",
			},
			checkErr: func(err error) bool {
				return errors.Is(err, errs.ErrInvalidPassword)
			},
			doNotMockUpdate: true,
		},
		{
			name: "not found user",
			req: models.UserChangePasswordReq{
				CurrentPassword: password,
				NewPassword:     newPassword,
				Login:           "testuser1",
			},
			mockGetErr: errs.ErrUserWithLoginNotFound,
			checkErr: func(err error) bool {
				return errors.Is(err, errs.ErrUserWithLoginNotFound)
			},
			doNotMockUpdate: true,
		},
		{
			name: "unexpected update error",
			req: models.UserChangePasswordReq{
				CurrentPassword: password,
				NewPassword:     newPassword,
				Login:           "testuser1",
			},
			mockUpdateErr: &pgconn.PgError{
				Code: pgerrcode.CannotConnectNow,
			},
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

			if !tt.doNotMock {
				repo.EXPECT().GetUserByLogin(context.Background(), tt.req.Login).
					Return(models.User{
						Login:        tt.req.Login,
						PasswordHash: hash,
					}, tt.mockGetErr)
			}

			if !tt.doNotMock && !tt.doNotMockUpdate {
				repo.EXPECT().UpdateUserPassword(context.Background(), mock.Anything).
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
