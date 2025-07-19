package middlewares

import (
	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	"github.com/LekcRg/GophKeeper/internal/server/repository"
	"go.uber.org/zap"
)

type Middlewares struct {
	log      *zap.Logger
	resp     *response.Responder
	config   *config.Config
	userRepo repository.UserRepo
}

func New(cfg *config.Config, log *zap.Logger, resp *response.Responder, userRepo repository.UserRepo) *Middlewares {
	return &Middlewares{
		log:      log,
		resp:     resp,
		config:   cfg,
		userRepo: userRepo,
	}
}
