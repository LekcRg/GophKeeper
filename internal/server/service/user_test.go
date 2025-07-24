package service

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/LekcRg/GophKeeper/internal/config"
// 	"github.com/LekcRg/GophKeeper/internal/crypto"
// 	"github.com/LekcRg/GophKeeper/internal/errs"
// 	"github.com/LekcRg/GophKeeper/internal/mocks"
// 	"github.com/LekcRg/GophKeeper/internal/models"
// 	validation "github.com/go-ozzo/ozzo-validation/v4"
// 	"github.com/jackc/pgerrcode"
// 	"github.com/jackc/pgx/v5/pgconn"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// type testUserServiceFunc func(ctx context.Context, req models.UserReq) (any, error)

// type testCase struct {
// 	mockErr   error
// 	checkErr  func(error) bool
// 	req       models.UserReq
// 	name      string
// 	doNotMock bool
// }

// func TestRegister(t *testing.T) {
// 	t.Parallel()

// 	cfg, err := config.GetConfig([]string{})
// 	require.NoError(t, err)

// 	tests := []testCase{
// 		{
// 			name: "success",
// 			req:  models.UserReq{Login: "testuser", Password: "T3stP@as5word"},
// 		},
// 		{
// 			name: "without password",
// 			req:  models.UserReq{Login: "testuser"},
// 			checkErr: func(err error) bool {
// 				var v validation.Errors
// 				return errors.As(err, &v)
// 			},
// 			doNotMock: true,
// 		},
// 		{
// 			name:    "login already exist",
// 			req:     models.UserReq{Login: "testuser", Password: "T3stP@as5word"},
// 			mockErr: errs.ErrLoginAlreadyExists,
// 			checkErr: func(err error) bool {
// 				return errors.Is(err, errs.ErrLoginAlreadyExists)
// 			},
// 		},
// 		{
// 			name:    "unexpected repository error",
// 			req:     models.UserReq{Login: "testuser", Password: "T3stP@as5word"},
// 			mockErr: &pgconn.PgError{Code: pgerrcode.CannotConnectNow},
// 			checkErr: func(err error) bool {
// 				var pgErr *pgconn.PgError
// 				return errors.As(err, &pgErr)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			runTestCase(t, cfg, &tt,
// 				func(repo *mocks.MockUserRepo) {
// 					repo.EXPECT().
// 						CreateUser(mock.Anything, mock.Anything).
// 						Return(42, tt.mockErr)
// 				},
// 				func(us *UserService) testUserServiceFunc {
// 					return func(ctx context.Context, req models.UserReq) (any, error) {
// 						return us.Register(ctx, req)
// 					}
// 				},
// 			)
// 		})
// 	}
// }

// func TestLogin(t *testing.T) {
// 	t.Parallel()

// 	cfg, err := config.GetConfig([]string{})
// 	require.NoError(t, err)

// 	const password = "T3stP@as5word"
// 	hash, err := crypto.HashPassword(password)
// 	require.NoError(t, err)

// 	tests := []testCase{
// 		{
// 			name: "success",
// 			req:  models.UserReq{Login: "testuser", Password: password},
// 		},
// 		{
// 			name: "without password",
// 			req:  models.UserReq{Login: "testuser"},
// 			checkErr: func(err error) bool {
// 				var v validation.Errors
// 				return errors.As(err, &v)
// 			},
// 			doNotMock: true,
// 		},
// 		{
// 			name: "invalid password",
// 			req:  models.UserReq{Login: "testuser", Password: "wrong"},
// 			checkErr: func(err error) bool {
// 				return errors.Is(err, errs.ErrInvalidCredentials)
// 			},
// 		},
// 		{
// 			name:    "not found",
// 			req:     models.UserReq{Login: "testuser", Password: password},
// 			mockErr: errs.ErrUserNotFound,
// 			checkErr: func(err error) bool {
// 				return errors.Is(err, errs.ErrInvalidCredentials)
// 			},
// 		},
// 		{
// 			name:    "unexpected error",
// 			req:     models.UserReq{Login: "testuser", Password: password},
// 			mockErr: &pgconn.PgError{Code: pgerrcode.CannotConnectNow},
// 			checkErr: func(err error) bool {
// 				var pgErr *pgconn.PgError
// 				return errors.As(err, &pgErr)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			runTestCase(t, cfg, &tt,
// 				func(repo *mocks.MockUserRepo) {
// 					repo.EXPECT().
// 						GetUserByLogin(mock.Anything, tt.req.Login).
// 						Return(models.User{
// 							ID:           42,
// 							Login:        tt.req.Login,
// 							PasswordHash: hash,
// 						}, tt.mockErr)
// 				},
// 				func(us *UserService) testUserServiceFunc {
// 					return func(ctx context.Context, req models.UserReq) (any, error) {
// 						return us.Login(ctx, req)
// 					}
// 				},
// 			)
// 		})
// 	}
// }

// func runTestCase(
// 	t *testing.T,
// 	cfg *config.Config,
// 	tt *testCase,
// 	mockFunc func(repo *mocks.MockUserRepo),
// 	svcFunc func(us *UserService) testUserServiceFunc,
// ) {
// 	t.Helper()

// 	repo := mocks.NewMockUserRepo(t)
// 	us := NewUserService(repo, cfg)

// 	if !tt.doNotMock {
// 		mockFunc(repo)
// 	}

// 	result, err := svcFunc(us)(context.Background(), tt.req)
// 	if tt.checkErr == nil {
// 		require.NoError(t, err)
// 		assert.NotNil(t, result)
// 	} else {
// 		assert.Error(t, err)
// 		assert.True(t, tt.checkErr(err))
// 	}
// }

// func TestChangePassword(t *testing.T) {
// 	t.Parallel()

// 	cfg, err := config.GetConfig([]string{})
// 	require.NoError(t, err)

// 	password := "T3stP@as5word"
// 	newPassword := "N3wSup3rP@55wrd"

// 	hash, err := crypto.HashPassword(password)
// 	require.NoError(t, err)

// 	type test struct {
// 		mockGetErr      error
// 		mockUpdateErr   error
// 		checkErr        func(error) bool
// 		req             models.UserChangePasswordReq
// 		name            string
// 		doNotMockGet    bool
// 		doNotMockUpdate bool
// 	}

// 	tests := []test{
// 		{
// 			name: "success",
// 			req: models.UserChangePasswordReq{
// 				Login:           "testuser",
// 				CurrentPassword: password,
// 				NewPassword:     newPassword,
// 			},
// 		},
// 		{
// 			name: "validation error",
// 			req: models.UserChangePasswordReq{
// 				Login:           "testuser",
// 				CurrentPassword: password,
// 				NewPassword:     "short",
// 			},
// 			checkErr: func(err error) bool {
// 				var v validation.Errors
// 				return errors.As(err, &v)
// 			},
// 			doNotMockGet: true,
// 		},
// 		{
// 			name: "invalid current password",
// 			req: models.UserChangePasswordReq{
// 				Login:           "testuser",
// 				CurrentPassword: "wrongPassword",
// 				NewPassword:     newPassword,
// 			},
// 			checkErr: func(err error) bool {
// 				return errors.Is(err, errs.ErrInvalidPassword)
// 			},
// 			doNotMockUpdate: true,
// 		},
// 		{
// 			name: "user not found",
// 			req: models.UserChangePasswordReq{
// 				Login:           "testuser",
// 				CurrentPassword: password,
// 				NewPassword:     newPassword,
// 			},
// 			mockGetErr: errs.ErrUserNotFound,
// 			checkErr: func(err error) bool {
// 				return errors.Is(err, errs.ErrUserNotFound)
// 			},
// 			doNotMockUpdate: true,
// 		},
// 		{
// 			name: "unexpected update error",
// 			req: models.UserChangePasswordReq{
// 				Login:           "testuser",
// 				CurrentPassword: password,
// 				NewPassword:     newPassword,
// 			},
// 			mockUpdateErr: &pgconn.PgError{Code: pgerrcode.CannotConnectNow},
// 			checkErr: func(err error) bool {
// 				var pgErr *pgconn.PgError
// 				return errors.As(err, &pgErr)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			repo := mocks.NewMockUserRepo(t)
// 			us := NewUserService(repo, cfg)

// 			if !tt.doNotMockGet {
// 				repo.EXPECT().
// 					GetUserByLogin(context.Background(), tt.req.Login).
// 					Return(models.User{
// 						ID:           1,
// 						Login:        tt.req.Login,
// 						PasswordHash: hash,
// 					}, tt.mockGetErr)
// 			}

// 			if !tt.doNotMockGet && !tt.doNotMockUpdate {
// 				repo.EXPECT().
// 					UpdateUserPassword(context.Background(), mock.Anything).
// 					Return(tt.mockUpdateErr)
// 			}

// 			err := us.ChangePassword(context.Background(), tt.req)
// 			if tt.checkErr == nil {
// 				require.NoError(t, err)
// 			} else {
// 				assert.Error(t, err)
// 				assert.True(t, tt.checkErr(err))
// 			}
// 		})
// 	}
// }
