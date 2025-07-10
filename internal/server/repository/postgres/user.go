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
	_ context.Context, _ models.CreateUserReq,
) (models.User, error) {
	return models.User{}, nil
}

type Msg struct {
	ID      int    `db:"id"`
	Message string `db:"test"`
}

func (ur *UserRepo) Test(ctx context.Context, message string) error {
	_, err := ur.db.ExecContext(ctx, "INSERT INTO test (test) VALUES($1);", message)
	if err != nil {
		return err
	}

	return nil
}
