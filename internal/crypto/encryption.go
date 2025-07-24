package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"runtime"

	"golang.org/x/crypto/argon2"
)

const (
	SaltLen    = 16
	TagContent = "__encrypted_tag__"
)

func GenEncryptionSalt() ([]byte, error) {
	salt := make([]byte, SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return salt, err
	}

	return salt, nil
}

func DeriveEncryptionKey(passwordStr string, salt []byte) []byte {
	const (
		time   uint32 = 5
		memory uint32 = 64 * 1024
		keyLen uint32 = 32
	)

	var (
		threads  = uint8(runtime.NumCPU())
		password = []byte(passwordStr)
	)

	return argon2.IDKey(password, salt, time, memory, threads, keyLen)
}

func Encrypt(content, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encrypted := aesGCM.Seal(nonce, nonce, content, nil)
	b64Encrypted := base64.StdEncoding.EncodeToString(encrypted)

	return b64Encrypted, nil
}

func Decrypt(encryptedString, password string, salt []byte) ([]byte, error) {
	key := DeriveEncryptionKey(password, salt)

	enc, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
