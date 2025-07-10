package service

import "github.com/LekcRg/GophKeeper/internal/server/repository"

type Service struct {
	UserService *UserService
}

func New(repo *repository.Repository) *Service {
	return &Service{
		UserService: NewUserService(repo.UserRepo),
	}
}
