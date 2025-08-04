package actions

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

func (a *Actions) CreateBinaryVault(ctx context.Context, name, path string) (models.VaultItem, error) {
	fileContent, fileInfo, err := a.readBinaryFile(path)
	if err != nil {
		return models.VaultItem{}, err
	}

	encMeta, err := a.encryptMetadata(fileInfo.Name(), fileInfo.Size())
	if err != nil {
		return models.VaultItem{}, err
	}

	req := models.VaultBinaryItemUploadReq{
		Name:          name,
		EncryptedData: encMeta,
	}

	res, err := a.req.CreateVaultBinaryItem(ctx, req)
	if err != nil {
		return models.VaultItem{}, err
	}

	if err := a.uploadEncryptedFile(ctx, res.URL, fileContent); err != nil {
		return models.VaultItem{}, err
	}

	if err := a.req.VaultConfirmCreateBinary(ctx, res.Path, res.ItemID); err != nil {
		return models.VaultItem{}, err
	}

	return res.Item, nil
}

func (a *Actions) readBinaryFile(path string) (content []byte, info os.FileInfo, err error) {
	var filePerm os.FileMode = 0o600

	file, err := os.OpenFile(path, os.O_RDONLY, filePerm)
	if err != nil {
		return nil, nil, err
	}

	defer file.Close()

	content, err = io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	info, err = file.Stat()
	if err != nil {
		return nil, nil, err
	}

	return content, info, nil
}

func (a *Actions) encryptMetadata(name string, size int64) (string, error) {
	data := models.VaultItemDataBinary{
		Path: name,
		Size: size,
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	enc, err := crypto.Encrypt(raw, a.state.CryptoKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(enc), nil
}

func (a *Actions) uploadEncryptedFile(ctx context.Context, url string, content []byte) error {
	enc, err := crypto.Encrypt(content, a.state.CryptoKey)
	if err != nil {
		return err
	}

	return a.req.VaultUploadBinaryFile(ctx, url, enc)
}

func (a *Actions) DownloadBinary(ctx context.Context, filename string, id int) (string, error) {
	file, err := a.DownloadBinaryBytes(ctx, id)
	if err != nil {
		return "", err
	}

	var (
		dirPerm  os.FileMode = 0o700
		filePerm os.FileMode = 0o600
	)

	path := filepath.Join(a.config.Folder, "binaries")

	err = os.MkdirAll(path, dirPerm)
	if err != nil {
		return "", fmt.Errorf("failed to create binaries dir: %w", err)
	}

	path = filepath.Join(path, filename)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, filePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write(file)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (a *Actions) DownloadBinaryBytes(ctx context.Context, id int) ([]byte, error) {
	url, err := a.req.VaultGetDowloadBidnaryURL(ctx, id)
	if err != nil {
		return nil, err
	}

	encFile, err := a.req.DownloadBinary(url)
	if err != nil {
		return nil, err
	}

	file, err := crypto.Decrypt(encFile, a.state.CryptoKey)
	if err != nil {
		return nil, err
	}

	return file, nil
}
