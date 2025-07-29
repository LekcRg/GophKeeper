package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"runtime"

	"github.com/LekcRg/GophKeeper/internal/errs"
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

func Decrypt(password string, enc, salt []byte) ([]byte, error) {
	key := DeriveEncryptionKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(enc) < nonceSize {
		return nil, errs.ErrInvalidEncrypted
	}

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func ValidEncryptionPassword(password string, tag, salt []byte) error {
	tag, err := Decrypt(password, tag, salt)
	if err != nil {
		return errs.ErrInvalidCryptoPasssword
	}

	tagStr := string(tag)
	if tagStr != TagContent {
		return errs.ErrInvalidCryptoPasssword
	}

	return nil
}
