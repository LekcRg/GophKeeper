package actions

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/models"
)

func (a *Actions) CreateVaultItem(
	ctx context.Context, name, vType string, data any,
) (models.VaultItem, error) {
	req := models.VaultCreateItemReq{
		Name: name,
		Type: vType,
	}

	jbytes, err := json.Marshal(data)
	if err != nil {
		return models.VaultItem{}, err
	}

	enc, err := crypto.Encrypt(jbytes, a.state.CryptoKey)
	if err != nil {
		return models.VaultItem{}, err
	}

	req.EncryptedData = base64.StdEncoding.EncodeToString(enc)

	return a.req.CreateVaultItem(ctx, req)
}
