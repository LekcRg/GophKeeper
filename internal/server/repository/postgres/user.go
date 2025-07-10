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

func (ur *UserRepo) GetUserByLogin(
	ctx context.Context, login string,
) (models.User, error) {
	query := "SELECT login, id, passhash FROM users WHERE login=$1"

	var user models.User

	err := ur.db.GetContext(ctx, &user, query, login)
	if err != nil {
		return user, err
	}

	return user, nil
}
