package actions

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

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

func (a *Actions) CreateBinaryVault(
	ctx context.Context, name, path string,
) (models.VaultItem, error) {
	var perm os.FileMode = 0o600

	file, err := os.OpenFile(path, os.O_RDONLY, perm)
	if err != nil {
		return models.VaultItem{}, err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return models.VaultItem{}, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return models.VaultItem{}, err
	}

	data := models.VaultItemDataBinary{
		Name: file.Name(),
		Size: fileInfo.Size(),
	}

	req := models.VaultBinaryItemUploadReq{
		Name: name,
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

	res, err := a.req.CreateVaultBinaryItem(ctx, req)
	if err != nil {
		return models.VaultItem{}, err
	}

	encFile, err := crypto.Encrypt(fileContent, a.state.CryptoKey)
	if err != nil {
		return models.VaultItem{}, err
	}

	err = a.req.VaultUploadBinaryFile(ctx, res.URL, encFile)
	if err != nil {
		return models.VaultItem{}, err
	}

	err = a.req.VaultConfirmCreateBinary(ctx, res.Path, res.ItemID)
	if err != nil {
		return models.VaultItem{}, err
	}

	return res.Item, nil
}
