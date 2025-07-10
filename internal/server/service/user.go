package service

import (
	"context"

	"github.com/LekcRg/GophKeeper/internal/server/repository"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(ur repository.UserRepo) *UserService {
	return &UserService{
		repo: ur,
	}
}

func (us *UserService) CreateUser(_ context.Context) error {
	return nil
}

func (us *UserService) Test(ctx context.Context, msg string) error {
	return us.repo.Test(ctx, msg)
}
