package postgres

import (
	"context"

	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) CreateUser(
	ctx context.Context, reqUser models.RegisterUserReq,
) error {
	query := "INSERT INTO users (login, passhash) VALUES (:login, :passhash)"

	_, err := ur.db.NamedExecContext(ctx, query, reqUser)
	if err != nil {
		return err
	}

	return nil
}
