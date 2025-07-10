package service

import (
	"context"
	"errors"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserService struct {
	repo   repository.UserRepo
	config *config.Config
}

func NewUserService(ur repository.UserRepo, cfg *config.Config) *UserService {
	return &UserService{
		repo:   ur,
		config: cfg,
	}
}

func (us *UserService) Register(ctx context.Context, req models.RegisterUserReq) (models.RegisterUserRes, error) {
	var (
		err error
		res models.RegisterUserRes
	)

	req.PasswordHash, err = crypto.HashPassword(req.Password)
	if err != nil {
		return res, err
	}

	err = us.repo.CreateUser(ctx, req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return res, errs.ErrLoginAlreadyExists
		}

		return res, err
	}

	res.Token, err = crypto.CreateJWTToken(req.Login, us.config.Auth)
	if err != nil {
		return res, err
	}

	return res, err
}
