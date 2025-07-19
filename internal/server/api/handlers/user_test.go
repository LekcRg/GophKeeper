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

			uh, svc := getHandlers(t)
			registerRunTest(t, tt, svc, uh.Register)
		})
	}
}

func registerRunTest(t *testing.T, tt testRegister, svc *MockUserService, h http.HandlerFunc) {
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

			uh, svc := getHandlers(t)
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

			if res.StatusCode() > 299 {
				compareErrs(t, tt.wantErrs, resErrs)
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

func getHandlers(t *testing.T) (*UserHandlers, *MockUserService) {
	t.Helper()

	svc := NewMockUserService(t)
	log := zaptest.NewLogger(t)
	resp := response.NewResponder(log)
	uh := NewUserHandlers(&config.Config{}, svc, log, resp)

	return uh, svc
}
