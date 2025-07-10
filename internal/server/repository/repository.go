package repository

import (
	"context"

	"github.com/LekcRg/GophKeeper/internal/models"
)

type DB interface {
	Close() error
}

type UserRepo interface {
	CreateUser(ctx context.Context, user models.RegisterUserReq) error
}

type Repository struct {
	DB       DB
	UserRepo UserRepo
}
