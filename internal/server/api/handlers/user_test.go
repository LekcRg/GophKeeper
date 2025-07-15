package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/api/middlewares"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"resty.dev/v3"
)

type testLoginRegister struct {
	mockSetup func(us *MockUserService)
	wantErrs  map[string]string
	name      string
	body      []byte
	wantCode  int
}

var userReq = models.UserReq{
	Login:    "testuser",
	Password: "testP@sw0rd123",
}

func TestRegister(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(userReq)
	require.NoError(t, err)

	tests := []testLoginRegister{
		{
			name: "success",
			body: body,
			mockSetup: func(svc *MockUserService) {
				svc.EXPECT().
					Register(mock.Anything, userReq).
					Return(models.TokenUserRes{Token: "token"}, nil)
			},
			wantCode: http.StatusCreated,
		},
		{
			name:     "invalid json",
			body:     []byte("{invalid json}"),
			wantCode: http.StatusBadRequest,
			wantErrs: map[string]string{"error": "Invalid JSON"},
		},
		{
			name: "internal service error",
			body: body,
			mockSetup: func(svc *MockUserService) {
				svc.EXPECT().
					Register(mock.Anything, userReq).
					Return(models.TokenUserRes{}, errors.New("Internal"))
			},
			wantCode: http.StatusInternalServerError,
			wantErrs: map[string]string{"error": "Internal server error"},
		},
		{
			name: "login already exists",
			body: body,
			mockSetup: func(svc *MockUserService) {
				svc.EXPECT().
					Register(mock.Anything, userReq).
					Return(models.TokenUserRes{}, errs.ErrLoginAlreadyExists)
			},
			wantCode: http.StatusConflict,
			wantErrs: map[string]string{"login": "Login already exists"},
		},
		{
			name: "validation error",
			body: body,
			mockSetup: func(svc *MockUserService) {
				svc.EXPECT().
					Register(mock.Anything, userReq).
					Return(models.TokenUserRes{}, validation.Errors{
						"password": errors.New("password requires special character"),
					})
			},
			wantCode: http.StatusBadRequest,
			wantErrs: map[string]string{"password": "password requires special character"},
		},
	}

	for _, ttest := range tests {
		tt := ttest
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uh, svc := getHandlers(t)
			loginRegisterRunTest(t, tt, svc, uh.Register)
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(userReq)
	require.NoError(t, err)

	tests := []testLoginRegister{
		{
			name: "success",
			body: body,
			mockSetup: func(svc *MockUserService) {
				svc.EXPECT().
					Login(mock.Anything, userReq).
					Return(models.TokenUserRes{Token: "token"}, nil)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "Invalid login or password",
			body: body,
			mockSetup: func(svc *MockUserService) {
				svc.EXPECT().
					Login(mock.Anything, userReq).
					Return(models.TokenUserRes{}, errs.ErrInvalidCredentials)
			},
			wantCode: http.StatusBadRequest,
			wantErrs: map[string]string{"login": "Invalid login or password"},
		},
	}

	for _, ttest := range tests {
		tt := ttest
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uh, svc := getHandlers(t)
			loginRegisterRunTest(t, tt, svc, uh.Login)
		})
	}
}

func loginRegisterRunTest(t *testing.T, tt testLoginRegister, svc *MockUserService, h http.HandlerFunc) {
	t.Helper()

	server := httptest.NewServer(h)
	defer server.Close()

	if tt.mockSetup != nil {
		tt.mockSetup(svc)
	}

	client := resty.New()
	defer client.Close()

	var (
		tokenRes = models.TokenUserRes{}
		resErrs  = map[string]string{}
		res      *resty.Response
		err      error
	)

	res, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(tt.body).
		SetResult(&tokenRes).
		SetError(&resErrs).
		Post(server.URL)

	assert.NoError(t, err)
	assert.Equal(t, tt.wantCode, res.StatusCode())

	if tt.wantCode < 300 {
		assert.NotEmpty(t, tokenRes.Token)
	} else {
		compareErrs(t, tt.wantErrs, resErrs)
	}
}

func TestChangePassword(t *testing.T) {
	t.Parallel()

	type test struct {
		mockSetup func(tc test, svc *MockUserService)
		wantErrs  map[string]string
		name      string
		login     string
		body      []byte
		wantCode  int
	}

	tests := []test{
		{
			name:     "success",
			login:    "testuser",
			body:     []byte(`{"current-password": "oldPass123!","new-passwrod":"newPass321!"}`),
			wantCode: http.StatusOK,
			mockSetup: func(tc test, svc *MockUserService) {
				var req models.UserChangePasswordReq
				err := json.Unmarshal(tc.body, &req)
				require.NoError(t, err)
				req.Login = tc.login

				svc.EXPECT().
					ChangePassword(mock.Anything, req).
					Return(nil)
			},
		},
		{
			name:     "validation error",
			login:    "testuser",
			body:     []byte(`{"current-password": "oldPass123!"}`),
			wantCode: http.StatusBadRequest,
			mockSetup: func(tc test, svc *MockUserService) {
				var req models.UserChangePasswordReq
				err := json.Unmarshal(tc.body, &req)
				require.NoError(t, err)
				req.Login = tc.login

				svc.EXPECT().
					ChangePassword(mock.Anything, req).
					Return(validation.Errors{
						"new-password": errors.New("new-password is requiered"),
					})
			},
			wantErrs: map[string]string{
				"new-password": "new-password is requiered",
			},
		},
		{
			name:     "invalid JSON",
			login:    "testuser",
			body:     []byte(`{"invalid json`),
			wantCode: http.StatusBadRequest,
			wantErrs: map[string]string{
				"error": "Invalid JSON",
			},
		},
		{
			name:     "unauthorized",
			body:     []byte(`{"current-password": "oldPass123!"}`),
			wantCode: http.StatusUnauthorized,
			wantErrs: map[string]string{
				"error": "Unauthorized",
			},
		},
		{
			name:     "service error",
			login:    "testuser",
			body:     []byte(`{"current-password": "oldPass123!"}`),
			wantCode: http.StatusBadRequest,
			mockSetup: func(tc test, svc *MockUserService) {
				var req models.UserChangePasswordReq
				err := json.Unmarshal(tc.body, &req)
				require.NoError(t, err)
				req.Login = tc.login

				svc.EXPECT().
					ChangePassword(mock.Anything, req).
					Return(errs.ErrInvalidPassword)
			},
			wantErrs: map[string]string{
				"current-password": "Password is not correct",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uh, svc := getHandlers(t)
			if tt.mockSetup != nil {
				tt.mockSetup(tt, svc)
			}

			server := httptest.NewServer(
				injectTestContext(
					http.HandlerFunc(uh.ChangePassword),
					tt.login),
			)
			defer server.Close()

			client := resty.New()
			defer client.Close()

			var (
				tokenRes = models.TokenUserRes{}
				resErrs  = map[string]string{}
			)

			res, _ := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(tt.body).
				SetResult(&tokenRes).
				SetError(&resErrs).
				Post(server.URL)

			assert.Equal(t, tt.wantCode, res.StatusCode())

			if res.StatusCode() > 299 {
				compareErrs(t, tt.wantErrs, resErrs)
			}
		})
	}
}

func injectTestContext(h http.Handler, login string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if login == "" {
			h.ServeHTTP(w, r)

			return
		}

		ctx := middlewares.AddLoginToCtx(r.Context(), login)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func compareErrs(t *testing.T, wantErrs, resErrs map[string]string) {
	for key, val := range wantErrs {
		assert.Equal(t, val, resErrs[key])
	}

	for key, val := range resErrs {
		assert.Equal(t, wantErrs[key], val)
	}
}

func getHandlers(t *testing.T) (*UserHandlers, *MockUserService) {
	t.Helper()

	svc := NewMockUserService(t)
	log := zaptest.NewLogger(t)
	resp := response.NewResponder(log)
	uh := NewUserHandlers(&config.Config{}, svc, log, resp)

	return uh, svc
}
