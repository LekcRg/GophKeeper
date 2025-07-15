package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"resty.dev/v3"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	type test struct {
		setupToken func() string
		name       string
		token      string
		wantCode   int
	}

	authCfg := config.Auth{Secret: "testsecret", JWTExpire: time.Minute}

	tests := []test{
		{
			name:     "success",
			wantCode: http.StatusOK,
			setupToken: func() string {
				token, err := crypto.CreateJWTToken("testuser", authCfg)
				require.NoError(t, err)

				return token
			},
		},
		{
			name:     "without token",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "invalid token",
			token:    "invalid",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "expired token",
			token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTI1ODI0MzMsImxvZ2luIjoidGVzdHVzZXIifQ.DS6RmGjqntbbXpIKE9Ikf0f_jEj4pys3rc-xEjNACkI",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "without login",
			wantCode: http.StatusUnauthorized,
			setupToken: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256,
					jwt.MapClaims{
						"exp": time.Now().Add(authCfg.JWTExpire).Unix(),
					})

				tokenString, err := token.SignedString([]byte(authCfg.Secret))
				require.NoError(t, err)

				return tokenString
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			token := tt.token
			if tt.setupToken != nil {
				token = "Bearer " + tt.setupToken()
			} else if token != "" {
				token = "Bearer " + token
			}

			log := zaptest.NewLogger(t)
			m := New(&config.Config{Auth: authCfg}, log, response.NewResponder(log))
			server := httptest.NewServer(m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("Success"))
				require.NoError(t, err)
			})))

			defer server.Close()

			t.Log(token)

			client := resty.New()
			defer client.Close()

			res, err := client.R().
				SetHeader("Authorization", token).
				Post(server.URL)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, res.StatusCode())
		})
	}
}

func TestLoginCtx(t *testing.T) {
	t.Parallel()
	t.Run("with login", func(t *testing.T) {
		t.Parallel()

		wantLogin := "testlogin"
		ctx := AddLoginToCtx(context.Background(), wantLogin)
		login, err := GetLogin(ctx)
		assert.NoError(t, err)
		assert.Equal(t, wantLogin, login)
	})

	t.Run("without login", func(t *testing.T) {
		t.Parallel()

		_, err := GetLogin(context.Background())
		assert.Error(t, err)
	})
}
