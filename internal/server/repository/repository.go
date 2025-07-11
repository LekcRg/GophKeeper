package repository

import (
	"context"

	"github.com/LekcRg/GophKeeper/internal/models"
)

type DB interface {
	Close() error
}

type UserRepo interface {
	CreateUser(ctx context.Context, user models.UserReq) error
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	UpdateUserPassword(ctx context.Context, user models.User) error
}

type Repository struct {
	DB       DB
	UserRepo UserRepo
}
