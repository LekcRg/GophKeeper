package main

import (
	"encoding/base64"
	"fmt"

	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/models"
	"resty.dev/v3"
)

func main() {
	const (
		login              = "testuser"
		password           = "I@mH3r00wow1337qq"
		encryptionPassword = "P@ssW0rd1234"
	)

	salt, err := crypto.GenEncryptionSalt()
	if err != nil {
		panic(err)
	}

	saltStr := base64.StdEncoding.EncodeToString(salt)
	fmt.Printf("salt: %s\n", saltStr)

	key := crypto.DeriveEncryptionKey(password, salt)

	enc, err := crypto.Encrypt([]byte(crypto.TagContent), key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("encrypted file base64: %s\n", enc)

	// dec, err := crypto.Decrypt(enc, password, salt)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("decrypted string: %s\n", dec)
	client := resty.New()
	response := models.APIKeyRes{}
	resErr := make(map[string]string, 0)

	res, err := client.R().SetBody(models.UserReq{
		Login:        login,
		Password:     password,
		EncryptedTag: enc,
		Salt:         saltStr,
	}).
		SetResult(&response).
		SetError(&resErr).
		Post("http://localhost:8080/user/create")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response status code: %d\n", res.StatusCode())
	fmt.Printf("Response: %+v\n", response)
	fmt.Printf("Response error: %+v\n", resErr)
}
