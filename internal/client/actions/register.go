package actions

import (
	"context"
	"encoding/base64"

	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/models"
)

func (a *Actions) Register(
	ctx context.Context, vals models.ClientAuthValues,
) (models.ClientRegisterResponse, error) {
	salt, err := crypto.GenEncryptionSalt()
	if err != nil {
		panic(err)
	}

	saltStr := base64.StdEncoding.EncodeToString(salt)
	key := crypto.DeriveEncryptionKey(vals.CryptoPassword, salt)

	encTag, err := crypto.Encrypt([]byte(crypto.TagContent), key)
	if err != nil {
		panic(err)
	}

	encTagString := base64.StdEncoding.EncodeToString(encTag)

	reqVals := models.UserReq{
		Salt:         saltStr,
		EncryptedTag: encTagString,
		Login:        vals.Login,
		Password:     vals.Password,
	}

	apiKey, err := a.req.UserRegister(ctx, reqVals)
	if err != nil {
		return models.ClientRegisterResponse{}, err
	}

	return models.ClientRegisterResponse{
		Key:  apiKey.Key,
		Salt: salt,
		Tag:  encTag,
	}, nil
}
