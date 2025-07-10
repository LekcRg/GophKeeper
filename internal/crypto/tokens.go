package crypto

import (
	"time"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWTToken(login string, cfg config.Auth) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"login": login,
			"exp":   time.Now().Add(cfg.JWTExpire).Unix(),
		})

	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
