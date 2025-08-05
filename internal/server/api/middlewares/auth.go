package middlewares

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/errs"
)

type key int

const (
	idKey key = iota
)

func (m *Middlewares) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		key = strings.TrimPrefix(key, "Bearer ")
		errorMsg := "Invalid token"

		if key == "" {
			m.resp.Error(w, http.StatusUnauthorized, errorMsg)
			return
		}

		var (
			splittedLen = 3
			splittedKey = strings.SplitN(key, "_", splittedLen)
		)

		if len(splittedKey) < splittedLen {
			m.resp.Error(w, http.StatusUnauthorized, errorMsg)
			return
		}

		id := splittedKey[1]
		hash := crypto.GenerateAPIHash(splittedKey[2])

		idInt, err := strconv.Atoi(id)
		if err != nil {
			m.resp.Error(w, http.StatusUnauthorized, errorMsg)
			return
		}

		user, err := m.userRepo.GetUserByID(r.Context(), idInt)
		if err != nil {
			m.resp.Error(w, http.StatusUnauthorized, errorMsg)
			return
		}

		if user.KeyHash != hash {
			m.resp.Error(w, http.StatusUnauthorized, errorMsg)
			return
		}

		ctx := AddIDToCtx(r.Context(), idInt)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

func AddIDToCtx(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, idKey, id)
}

func GetID(ctx context.Context) (int, error) {
	id, ok := ctx.Value(idKey).(int)
	if !ok {
		return 0, errs.ErrNotValidContextID
	}

	return id, nil
}
