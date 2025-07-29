package actions

import (
	"context"
	"encoding/base64"

	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/models"
	"go.uber.org/zap"
)

func (a *Actions) Register(
	ctx context.Context, vals models.ClientRegisterValues,
) (models.ClientRegisterResponse, error) {
	salt, err := crypto.GenEncryptionSalt()
	if err != nil {
		a.log.Error("error generating salt", zap.Error(err))
		return models.ClientRegisterResponse{}, err
	}

	saltStr := base64.StdEncoding.EncodeToString(salt)
	key := crypto.DeriveEncryptionKey(vals.CryptoPassword, salt)

	encTag, err := crypto.Encrypt([]byte(crypto.TagContent), key)
	if err != nil {
		a.log.Error("error generating encrypted tag", zap.Error(err))
		return models.ClientRegisterResponse{}, err
	}

	encTagString := base64.StdEncoding.EncodeToString(encTag)

	reqVals := models.UserReq{
		Salt:         saltStr,
		EncryptedTag: encTagString,
		Login:        vals.Login,
		Password:     vals.Password,
	}

	apiKey, err := a.req.Register(ctx, reqVals)
	if err != nil {
		a.log.Error("User register request error", zap.Error(err))
		return models.ClientRegisterResponse{}, err
	}

	return models.ClientRegisterResponse{
		Key:  apiKey.Key,
		Salt: salt,
		Tag:  encTag,
	}, nil
}

func (a *Actions) UpdateKey(ctx context.Context, req models.UserLogin) (models.APIKeyRes, error) {
	return a.req.UpdateAPIKey(ctx, req)
}
