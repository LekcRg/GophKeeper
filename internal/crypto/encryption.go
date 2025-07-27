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
		numCPU         = runtime.NumCPU()
		threads  uint8 = 4
		password       = []byte(passwordStr)
	)

	if numCPU > 0 && numCPU <= 255 {
		threads = uint8(numCPU)
	}

	return argon2.IDKey(password, salt, time, memory, threads, keyLen)
}

func Encrypt(content, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrypted := aesGCM.Seal(nonce, nonce, content, nil)

	return encrypted, nil
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
