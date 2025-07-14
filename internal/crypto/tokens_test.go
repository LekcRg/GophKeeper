package crypto

import (
	"strings"
	"testing"
	"time"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTToken(t *testing.T) {
	t.Parallel()

	type test struct {
		name  string
		login string
		cfg   config.Auth
	}

	tests := []test{
		{
			name:  "simple login and secret",
			login: "test",
			cfg: config.Auth{
				Secret:    "secret",
				JWTExpire: 1 * time.Hour,
			},
		},
		{
			name:  "long login",
			login: strings.Repeat("s", 1000),
			cfg: config.Auth{
				Secret:    "secret",
				JWTExpire: 1 * time.Hour,
			},
		},
		{
			name:  "emoji in login",
			login: "emoji🙄",
			cfg: config.Auth{
				Secret:    "secret",
				JWTExpire: 1 * time.Hour,
			},
		},
	}

	for _, ttest := range tests {
		tt := ttest
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			exp := time.Now().Add(tt.cfg.JWTExpire - time.Minute)
			token, err := CreateJWTToken(tt.login, tt.cfg)
			require.NoError(t, err)
			assert.NotEmpty(t, token)

			claims, err := ValidJWTToken(token, tt.cfg.Secret)
			require.NoError(t, err)

			jwtExp, err := claims.GetExpirationTime()
			require.NoError(t, err)
			assert.Greater(t, jwtExp.Time, exp)

			jwtLogin, ok := claims["login"]
			require.True(t, ok)
			assert.Equal(t, tt.login, jwtLogin)
		})
	}
}
