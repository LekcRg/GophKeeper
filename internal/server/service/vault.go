package service

import (
	"context"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/repository"
	"github.com/LekcRg/GophKeeper/internal/server/service/valid"
)

type VaultService struct {
	repo   repository.VaultRepo
	config *config.Config
}

func NewVaultService(vr repository.VaultRepo, cfg *config.Config) *VaultService {
	return &VaultService{
		repo:   vr,
		config: cfg,
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
	return vs.repo.GetAllItems(ctx, id)
}
