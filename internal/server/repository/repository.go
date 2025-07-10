package repository

import (
	"context"

	"github.com/LekcRg/GophKeeper/internal/models"
)

type DB interface {
	Close() error
}

type UserRepo interface {
	CreateUser(ctx context.Context, user models.CreateUserReq) (models.User, error)
	Test(ctx context.Context, message string) error
}

type Repository struct {
	DB       DB
	UserRepo UserRepo
}
