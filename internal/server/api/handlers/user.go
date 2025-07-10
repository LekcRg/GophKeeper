package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/LekcRg/GophKeeper/internal/config"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(ctx context.Context) error
}

type UserHandlers struct {
	service UserService
	config  *config.Config
	log     *zap.Logger
}

func NewUserHandlers(cfg *config.Config, service UserService, log *zap.Logger) *UserHandlers {
	return &UserHandlers{
		service: service,
		config:  cfg,
		log:     log,
	}
}

func (us *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	us.log.Info("CreteUser handler", zap.String("URI", r.RequestURI))

	w.WriteHeader(http.StatusOK)

	_, err := io.WriteString(w, "Hello from user handler")
	if err != nil {
		us.log.Error("CreateUser error", zap.Error(err))
	}
}
