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
) (models.APIKeyRes, error) {
	var (
		err error
		res models.APIKeyRes
	)

	err = valid.Register(&req)
	if err != nil {
		return res, err
	}

	req.PasswordHash, err = crypto.HashPassword(req.Password)
	if err != nil {
		return res, err
	}

	random, hash, err := crypto.CreateRandomPartAPIKey(us.config.Auth)
	if err != nil {
		return res, err
	}

	req.KeyHash = hash

	id, err := us.repo.CreateUser(ctx, req)
	if err != nil {
		return res, err
	}

	res.Key = crypto.JoinFullAPIKey(id, random)

	return res, err
}

func (us *UserService) Login(
	ctx context.Context, req models.UserReq,
) (int, error) {
	err := valid.Login(&req)
	if err != nil {
		return 0, err
	}

	user, err := us.repo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return 0, errs.ErrInvalidCredentials
		}

		return 0, err
	}

	isValid := crypto.CheckPasswordHash(req.Password, user.PasswordHash)
	if !isValid {
		return 0, errs.ErrInvalidCredentials
	}

	return user.ID, err
}

func (us *UserService) UpdateAPIKey(ctx context.Context, req models.UserReq) (models.APIKeyRes, error) {
	res := models.APIKeyRes{}

	id, err := us.Login(ctx, req)
	if err != nil {
		return res, err
	}

	var hash string

	res.Key, hash, err = crypto.CreateFullAPIKey(id, us.config.Auth)
	if err != nil {
		return res, err
	}

	err = us.repo.UpdateUserKey(ctx, models.User{
		ID:      id,
		KeyHash: hash,
	})
	if err != nil {
		return res, err
	}

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
