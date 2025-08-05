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
) (int, error) {
	query := `INSERT INTO users (login, passhash, key_hash, encrypted_tag, salt)
	VALUES (:login, :passhash, :key_hash, :encrypted_tag, :salt) RETURNING id`

	rows, err := ur.db.NamedQueryContext(ctx, query, reqUser)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, errs.ErrLoginAlreadyExists
		}

		return 0, err
	}

	defer rows.Close()

	var id int
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (ur *UserRepo) GetUserByLogin(
	ctx context.Context, login string,
) (models.User, error) {
	query := "SELECT * FROM users WHERE login=$1"

	var user models.User

	err := ur.db.GetContext(ctx, &user, query, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errs.ErrUserNotFound
		}

		return user, err
	}

	return user, nil
}

func (ur *UserRepo) GetUserByID(
	ctx context.Context, id int,
) (models.User, error) {
	query := "SELECT * FROM users WHERE id=$1"

	var user models.User

	err := ur.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errs.ErrUserNotFound
		}

		return user, err
	}

	return user, nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, query string, user models.User) error {
	_, err := ur.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) UpdateUserKey(ctx context.Context, user models.User) error {
	query := "UPDATE users SET key_hash = :key_hash WHERE id = :id"

	return ur.UpdateUser(ctx, query, user)
}

func (ur *UserRepo) UpdateUserPassword(
	ctx context.Context, user models.User,
) error {
	query := "UPDATE users SET passhash = :passhash WHERE login = :login"

	return ur.UpdateUser(ctx, query, user)
}
