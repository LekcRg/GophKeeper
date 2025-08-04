package handlers

import (
	"context"
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
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"resty.dev/v3"
)

type testRegister struct {
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

	tests := []testRegister{
		{
			name: "success",
			body: body,
			mockSetup: func(svc *MockUserService) {
				svc.EXPECT().
					Register(mock.Anything, userReq).
					Return(models.APIKeyRes{Key: "key"}, nil)
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
					Return(models.APIKeyRes{}, errors.New("Internal"))
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
					Return(models.APIKeyRes{}, errs.ErrLoginAlreadyExists)
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
					Return(models.APIKeyRes{}, validation.Errors{
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

			uh, svc := getUserHandlers(t)
			userReqRunTest(t, tt, svc, uh.Register)
		})
	}
}

func TestAPIKey(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(userReq)
	require.NoError(t, err)

	data := models.UserLogin{
		Login:    "testuser",
		Password: "testP@sw0rd123",
	}

	tests := []testRegister{
		{
			name: "success",
			body: body,
			mockSetup: func(svc *MockUserService) {
				svc.EXPECT().
					UpdateAPIKey(mock.Anything, data).
					Return(models.APIKeyRes{Key: "key"}, nil)
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
					UpdateAPIKey(mock.Anything, data).
					Return(models.APIKeyRes{}, errors.New("Internal"))
			},
			wantCode: http.StatusInternalServerError,
			wantErrs: map[string]string{"error": "Internal server error"},
		},
		{
			name: "invalid login or password",
			body: body,
			mockSetup: func(svc *MockUserService) {
				svc.EXPECT().
					UpdateAPIKey(mock.Anything, data).
					Return(models.APIKeyRes{}, errs.ErrInvalidCredentials)
			},
			wantCode: http.StatusBadRequest,
			wantErrs: map[string]string{"login": "Invalid login or password"},
		},
	}

	for _, ttest := range tests {
		tt := ttest
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uh, svc := getUserHandlers(t)
			userReqRunTest(t, tt, svc, uh.APIKey)
		})
	}
}

func userReqRunTest(t *testing.T, tt testRegister, svc *MockUserService, h http.HandlerFunc) {
	t.Helper()

	server := httptest.NewServer(h)
	defer server.Close()

	if tt.mockSetup != nil {
		tt.mockSetup(svc)
	}

	client := resty.New()
	defer client.Close()

	var (
		tokenRes = models.APIKeyRes{}
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
		assert.NotEmpty(t, tokenRes.Key)
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
			body:     []byte(`{"login": "testuser", "current-password": "oldPass123!","new-passwrod":"newPass321!"}`),
			wantCode: http.StatusOK,
			mockSetup: func(tc test, svc *MockUserService) {
				var req models.UserChangePasswordReq
				err := json.Unmarshal(tc.body, &req)
				require.NoError(t, err)

				svc.EXPECT().
					ChangePassword(mock.Anything, req).
					Return(nil)
			},
		},
		{
			name:     "validation error",
			login:    "testuser",
			body:     []byte(`{"login": "testuser", "current-password": "oldPass123!"}`),
			wantCode: http.StatusBadRequest,
			mockSetup: func(tc test, svc *MockUserService) {
				var req models.UserChangePasswordReq
				err := json.Unmarshal(tc.body, &req)
				require.NoError(t, err)

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
			name:     "service error",
			login:    "testuser",
			body:     []byte(`{"login": "testuser", "current-password": "oldPass123!"}`),
			wantCode: http.StatusBadRequest,
			mockSetup: func(tc test, svc *MockUserService) {
				var req models.UserChangePasswordReq
				err := json.Unmarshal(tc.body, &req)
				require.NoError(t, err)

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

			uh, svc := getUserHandlers(t)
			if tt.mockSetup != nil {
				tt.mockSetup(tt, svc)
			}

			server := httptest.NewServer(
				http.HandlerFunc(uh.ChangePassword),
			)
			defer server.Close()

			client := resty.New()
			defer client.Close()

			var (
				tokenRes = models.Response{}
				resErrs  = map[string]string{}
			)

			res, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(tt.body).
				SetResult(&tokenRes).
				SetError(&resErrs).
				Post(server.URL)

			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, res.StatusCode())

			if !res.IsSuccess() {
				compareErrs(t, tt.wantErrs, resErrs)
			}
		})
	}
}

func TestGetCryproParams(t *testing.T) {
	t.Parallel()

	type test struct {
		ctx      context.Context
		svcErr   error
		wantErrs map[string]string
		svcRes   models.CryptoParamsRes
		name     string
		wantRes  string
		id       int
		wantCode int
		svcMock  bool
	}

	tests := []test{
		{
			name:     "Success",
			ctx:      middlewares.AddIDToCtx(context.Background(), 1),
			id:       1,
			wantCode: http.StatusOK,
			svcRes: models.CryptoParamsRes{
				EncryptedTag: "tag",
				Salt:         "salt",
			},
			svcMock: true,
			wantRes: `{"encrypted_tag":"tag","salt":"salt"}`,
		},
		{
			name:     "Without ID",
			ctx:      context.Background(),
			wantCode: http.StatusUnauthorized,
			wantErrs: map[string]string{"error": "Unauthorized"},
		},
		{
			name:     "Internal server error",
			ctx:      middlewares.AddIDToCtx(context.Background(), 1),
			id:       1,
			wantCode: http.StatusInternalServerError,
			svcErr:   errors.New("error"),
			svcMock:  true,
			wantErrs: map[string]string{"error": "Internal server error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uh, svc := getUserHandlers(t)

			if tt.svcMock {
				svc.EXPECT().GetCryptoParams(tt.ctx, tt.id).Return(tt.svcRes, tt.svcErr)
			}

			res := serveHTTPWithCtx(tt.ctx, uh.GetCryptoParams, serveHTTPOpts{
				method: "GET",
			})

			assert.Equal(t, tt.wantCode, res.Code)

			if res.Code > 299 {
				var resErrs map[string]string
				err := json.Unmarshal(res.Body.Bytes(), &resErrs)
				require.NoError(t, err)
				compareErrs(t, tt.wantErrs, resErrs)
			} else {
				assert.Equal(t, tt.wantRes, res.Body.String())
			}
		})
	}
}

func compareErrs(t *testing.T, wantErrs, resErrs map[string]string) {
	t.Helper()

	for key, val := range wantErrs {
		assert.Equal(t, val, resErrs[key])
	}

	for key, val := range resErrs {
		assert.Equal(t, wantErrs[key], val)
	}
}

func getUserHandlers(t *testing.T) (*UserHandlers, *MockUserService) {
	t.Helper()

	svc := NewMockUserService(t)
	log := zaptest.NewLogger(t)
	resp := response.NewResponder(log)
	uh := NewUserHandlers(&config.Config{}, svc, log, resp)

	return uh, svc
}
