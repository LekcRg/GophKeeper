package handlers

import (
	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	"github.com/LekcRg/GophKeeper/internal/server/service"
	"go.uber.org/zap"
)

type Handlers struct {
	UserHandlers *UserHandlers
}

func New(cfg *config.Config, svc *service.Service, log *zap.Logger, resp *response.Responder) *Handlers {
	return &Handlers{
		UserHandlers: NewUserHandlers(cfg, svc.UserService, log, resp),
	}
}
