package repository

import (
	"context"
	"database/sql"

	"github.com/LekcRg/GophKeeper/internal/models"
)

type DB interface {
	Close() error
}

type UserRepo interface {
	CreateUser(ctx context.Context, user models.UserReq) (int, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	UpdateUserKey(ctx context.Context, user models.User) error
	UpdateUserPassword(ctx context.Context, user models.User) error
}

type Repository struct {
	DB       DB
	SQL      *sql.DB
	UserRepo UserRepo
}
