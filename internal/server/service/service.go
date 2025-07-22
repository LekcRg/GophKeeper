package service

import (
	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/server/repository"
)

type Service struct {
	UserService  *UserService
	VaultService *VaultService
}

func New(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		UserService:  NewUserService(repo.UserRepo, cfg),
		VaultService: NewVaultService(repo.VaultRepo, cfg),
	}
}
