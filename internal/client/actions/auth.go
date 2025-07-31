package actions

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/req"
	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

func (a *Actions) Register(
	ctx context.Context, vals models.ClientRegisterValues,
) (msgs.CredentialsBytesMsg, error) {
	salt, err := crypto.GenEncryptionSalt()
	if err != nil {
		a.log.Error("error generating salt", zap.Error(err))
		return msgs.CredentialsBytesMsg{}, err
	}

	saltStr := base64.StdEncoding.EncodeToString(salt)
	key := crypto.DeriveEncryptionKey(vals.CryptoPassword, salt)

	encTag, err := crypto.Encrypt([]byte(crypto.TagContent), key)
	if err != nil {
		a.log.Error("error generating encrypted tag", zap.Error(err))
		return msgs.CredentialsBytesMsg{}, err
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
		return msgs.CredentialsBytesMsg{}, err
	}

	return msgs.CredentialsBytesMsg{
		Key:  apiKey.Key,
		Salt: salt,
		Tag:  encTag,
	}, nil
}

func (a *Actions) GetCredentials(
	ctx context.Context, key string,
) (msgs.CredentialsBytesMsg, error) {
	res, err := a.req.GetCredentials(ctx, key)
	if err != nil {
		return msgs.CredentialsBytesMsg{}, err
	}

	saltBytes, err := base64.StdEncoding.DecodeString(res.Salt)
	if err != nil {
		return msgs.CredentialsBytesMsg{}, err
	}

	tagBytes, err := base64.StdEncoding.DecodeString(res.EncryptedTag)
	if err != nil {
		return msgs.CredentialsBytesMsg{}, err
	}

	return msgs.CredentialsBytesMsg{
		Key:  key,
		Salt: saltBytes,
		Tag:  tagBytes,
	}, nil
}

func (a *Actions) UpdateKey(ctx context.Context, request models.UserLogin) (models.APIKeyRes, error) {
	return a.req.UpdateAPIKey(ctx, request)
}

func (a *Actions) CheckCryptoPassword(form msgs.FormSubmitMsg, passwordInput string) tea.Cmd {
	return func() tea.Msg {
		err := crypto.ValidEncryptionPassword(form.Values[passwordInput], a.config.EnctyptedTag, a.config.Salt)
		if err != nil {
			if errors.Is(err, errs.ErrInvalidCryptoPasssword) {
				return &req.ResError{
					Errors: map[string]string{
						passwordInput: "invalid password",
					},
				}
			}

			return msgs.ErrorMsg(err)
		}

		return msgs.CryptoPassValid(passwordInput)
	}
}
