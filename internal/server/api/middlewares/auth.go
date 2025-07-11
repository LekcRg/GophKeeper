package middlewares

import (
	"context"
	"net/http"

	"github.com/LekcRg/GophKeeper/internal/crypto"
	"go.uber.org/zap"
)

func (m *Middlewares) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		m.log.Info(token)
		if token == "" {
			m.resp.Error(w, http.StatusUnauthorized, "Unauthorized")

			return
		}

		claim, err := crypto.ValidJWTToken(token, m.config.Auth.Secret)
		if err != nil {
			m.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
			m.log.Info("Invalid JWT", zap.Error(err))

			return
		}

		login, ok := claim["login"]
		if !ok {
			m.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
			m.log.Error("Login not found in JWT token")

			return
		}

		loginStr, ok := login.(string)
		if !ok {
			m.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
			m.log.Error("Login not string in JWT token")

			return
		}

		ctx := context.WithValue(r.Context(), "login", loginStr)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
