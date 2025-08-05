package crypto

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 10
)

var (
	letters  = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	specials = []rune("!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~")
	nums     = []rune("0123456789")
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func GeneratePassword(length int, useLetters, useSpecial bool) (string, error) {
	alphabet := nums
	if useLetters {
		alphabet = append(alphabet, letters...)
	}

	if useSpecial {
		alphabet = append(alphabet, specials...)
	}

	res := make([]rune, length)

	for i := range length {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}

		res[i] = alphabet[r.Int64()]
	}

	return string(res), nil
}
