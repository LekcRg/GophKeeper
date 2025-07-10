package service

import "context"

type UserService struct {
	// test string
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) CreateUser(_ context.Context) error {
	return nil
}
