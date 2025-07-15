package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
	ctx context.Context, reqUser models.UserReq,
) error {
	query := "INSERT INTO users (login, passhash) VALUES (:login, :passhash)"

	_, err := ur.db.NamedExecContext(ctx, query, reqUser)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errs.ErrLoginAlreadyExists
		}

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
		if errors.Is(err, sql.ErrNoRows) {
			return user, errs.ErrUserWithLoginNotFound
		}

		return user, err
	}

	return user, nil
}

func (ur *UserRepo) UpdateUserPassword(
	ctx context.Context, user models.User,
) error {
	query := "UPDATE users SET passhash = :passhash WHERE login = :login"

	_, err := ur.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	return nil
}
