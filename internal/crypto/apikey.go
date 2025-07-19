package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/LekcRg/GophKeeper/internal/config"
)

func GenerateAPIHash(key string) string {
	hash := sha256.Sum256([]byte(key))

	return fmt.Sprintf("%x", hash)
}

func CreateRandomPartAPIKey(
	cfg config.Auth,
) (
	encoded, hash string, err error,
) {
	randomBytes := make([]byte, cfg.MaxBytesKey)

	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", "", err
	}

	encoded = base64.RawURLEncoding.EncodeToString(randomBytes)
	hash = GenerateAPIHash(encoded)

	return encoded, hash, nil
}

func CreateFullAPIKey(
	id int, cfg config.Auth,
) (
	key, hash string, err error,
) {
	random, hash, err := CreateRandomPartAPIKey(cfg)
	if err != nil {
		return "", "", err
	}

	token := JoinFullAPIKey(id, random)

	return token, hash, nil
}

func JoinFullAPIKey(id int, random string) string {
	return fmt.Sprintf("sk_%d_%s", id, random)
}
