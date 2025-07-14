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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"resty.dev/v3"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	type test struct {
		mockSetup func(us *MockUserService)
		wantErrs  map[string]string
		name      string
		body      []byte
		wantCode  int
	}

	req := models.UserReq{
		Login:    "testuser",
		Password: "testP@sw0rd123",
	}
	body, err := json.Marshal(req)
	require.NoError(t, err)

	tests := []test{
		{
			name: "success",
			body: body,
			mockSetup: func(svc *MockUserService) {
				svc.EXPECT().
					Register(mock.Anything, req).
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
					Register(mock.Anything, req).
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
					Register(mock.Anything, req).
					Return(models.TokenUserRes{}, errs.ErrLoginAlreadyExists)
			},
			wantCode: http.StatusConflict,
			wantErrs: map[string]string{"login": "login already exists"},
		},
		{
			name: "validation error",
			body: body,
			mockSetup: func(svc *MockUserService) {
				svc.EXPECT().
					Register(mock.Anything, req).
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

			svc := NewMockUserService(t)
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}

			log := zaptest.NewLogger(t)
			resp := response.NewResponder(log)
			uh := NewUserHandlers(&config.Config{}, svc, log, resp)

			server := httptest.NewServer(http.HandlerFunc(uh.Register))
			defer server.Close()

			client := resty.New()
			defer client.Close()

			var (
				tokenRes = models.TokenUserRes{}
				ers      = map[string]string{}
				res      *resty.Response
				err      error
			)

			res, err = client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(tt.body).
				SetResult(&tokenRes).
				SetError(&ers).
				Post(server.URL)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCode, res.StatusCode())

			if tt.wantCode < 300 {
				assert.NotEmpty(t, tokenRes.Token)
			} else {
				for key, e := range tt.wantErrs {
					assert.Equal(t, e, ers[key])
				}
			}
		})
	}
}
