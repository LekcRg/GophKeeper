package service

import (
	"context"
	"errors"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/repository"
	"github.com/LekcRg/GophKeeper/internal/server/service/valid"
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

func (us *UserService) Register(
	ctx context.Context, req models.UserReq,
) (models.TokenUserRes, error) {
	var (
		err error
		res models.TokenUserRes
	)

	err = valid.Register(&req)
	if err != nil {
		return res, err
	}

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

	return res, err
}

func (us *UserService) Login(
	ctx context.Context, req models.UserReq,
) (models.TokenUserRes, error) {
	res := models.TokenUserRes{}

	user, err := us.repo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return res, err
	}

	isValid := crypto.CheckPasswordHash(req.Password, user.PasswordHash)
	if !isValid {
		return res, errs.ErrInvalidCredentials
	}

	res.Token, err = crypto.CreateJWTToken(user.Login, us.config.Auth)

	return res, err
}

func (us *UserService) ChangePassword(ctx context.Context, req models.UserChangePasswordReq) error {
	err := valid.ChangePassword(&req)
	if err != nil {
		return err
	}

	user, err := us.repo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return err
	}

	validPass := crypto.CheckPasswordHash(req.CurrentPassword, user.PasswordHash)
	if !validPass {
		return errs.ErrInvalidPassword
	}

	newPassHash, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newPassHash

	err = us.repo.UpdateUserPassword(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
