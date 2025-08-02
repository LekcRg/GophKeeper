package service

import (
	"context"
	"encoding/base64"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/repository"
	"github.com/LekcRg/GophKeeper/internal/server/service/valid"
	"github.com/LekcRg/GophKeeper/internal/server/storage"
)

type VaultService struct {
	repo    repository.VaultRepo
	config  *config.Config
	storage *storage.Storage
}

func NewVaultService(
	vr repository.VaultRepo, cfg *config.Config, st *storage.Storage,
) *VaultService {
	return &VaultService{
		repo:    vr,
		config:  cfg,
		storage: st,
	}
}

func (vs *VaultService) CreateItem(
	ctx context.Context, item models.VaultItem,
) (models.VaultItem, error) {
	err := valid.VaultCreateItem(&item)
	if err != nil {
		return models.VaultItem{}, err
	}

	res, err := vs.repo.CreateItem(ctx, item)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (vs *VaultService) GetAllItems(
	ctx context.Context, id int,
) ([]models.VaultItem, error) {
	res, err := vs.repo.GetAllItems(ctx, id)
	if err != nil {
		return []models.VaultItem{}, err
	}

	for i := range res {
		item := &res[i]
		item.EncryptedDataString = base64.StdEncoding.EncodeToString(item.EncryptedData)
	}

	return res, nil
}

func (vs *VaultService) CreateBinary(ctx context.Context, item models.VaultItem) (
	models.VaultBinaryItemUploadRes, error,
) {
	createdItem, err := vs.CreateItem(ctx, item)
	if err != nil {
		return models.VaultBinaryItemUploadRes{}, err
	}

	url, path, err := vs.storage.GenUploadPresignedUrl(ctx, item.UserID)
	if err != nil {
		return models.VaultBinaryItemUploadRes{}, err
	}

	return models.VaultBinaryItemUploadRes{
		Item:   createdItem,
		ItemID: createdItem.ID,
		URL:    url,
		Path:   path,
	}, nil
}

func (vs *VaultService) ConfirmBinaryUpload(
	ctx context.Context, req models.VaultConfirmBinaryUploadReq,
) error {
	err := valid.ValidConfirmBinaryUpload(&req)
	if err != nil {
		return err
	}

	if !vs.storage.IsContainsFile(ctx, req.Path) {
		return errs.ErrBinaryFileNotFound
	}

	return vs.repo.UpdateBinaryURL(ctx, req)
}

func (vs *VaultService) GetBinaryFileURL(
	ctx context.Context, userID, vaultID int,
) (string, error) {
	item, err := vs.repo.GetItem(ctx, vaultID)
	if err != nil {
		return "", err
	}

	if item.UserID != userID {
		return "", errs.ErrInvalidUserBinary
	}

	url, err := vs.storage.GenPresignedGetUrl(ctx, item.BinaryPath)
	if err != nil {
		return "", err
	}

	return url, nil
}
