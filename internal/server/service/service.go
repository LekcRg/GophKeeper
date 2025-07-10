package service

type Service struct {
	UserService *UserService
}

func New() *Service {
	return &Service{
		UserService: NewUserService(),
	}
}
