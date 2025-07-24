package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/mocks"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"resty.dev/v3"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	type test struct {
		mockErr  error
		name     string
		key      string
		wantCode int
		doMock   bool
	}

	cfg, err := config.GetConfig([]string{})
	require.NoError(t, err)

	key, hash, err := crypto.CreateFullAPIKey(1, cfg.Auth)
	require.NoError(t, err)

	tests := []test{
		{
			name:     "success",
			wantCode: http.StatusOK,
			key:      "Bearer " + key,
			doMock:   true,
		},
		{
			name:     "without key",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "1 invalid key",
			key:      "invalid",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "2 invalid key",
			key:      "sk_1_asdlfkjaskdjfasdf",
			wantCode: http.StatusUnauthorized,
			doMock:   true,
		},
		{
			name:     "3 invalid key",
			key:      "sk_asdfasf_asdlfkjaskdjfasdf",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "Not found user",
			key:      "Bearer " + key,
			doMock:   true,
			mockErr:  errs.ErrUserNotFound,
			wantCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			log := zaptest.NewLogger(t)
			repo := mocks.NewMockUserRepo(t)

			if tt.doMock {
				repo.EXPECT().GetUserByID(mock.Anything, 1).Return(models.User{
					ID:           1,
					Login:        "testuser",
					PasswordHash: "asdf",
					KeyHash:      hash,
				}, tt.mockErr)
			}

			m := New(cfg, log, response.NewResponder(log), repo)
			server := httptest.NewServer(m.Authenticate(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte("Success"))
					require.NoError(t, err)
				})))

			defer server.Close()

			client := resty.New()
			defer client.Close()

			res, err := client.R().
				SetHeader("Authorization", tt.key).
				Post(server.URL)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, res.StatusCode())
		})
	}
}

func TestIDCtx(t *testing.T) {
	t.Parallel()
	t.Run("with ID", func(t *testing.T) {
		t.Parallel()

		wantID := 23
		ctx := AddIDToCtx(context.Background(), wantID)
		login, err := GetID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, wantID, login)
	})

	t.Run("without ID", func(t *testing.T) {
		t.Parallel()

		_, err := GetID(context.Background())
		assert.Error(t, err)
	})
}
