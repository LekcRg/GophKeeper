package middlewares

import (
	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	"go.uber.org/zap"
)

type Middlewares struct {
	log    *zap.Logger
	resp   *response.Responder
	config *config.Config
}

func New(cfg *config.Config, log *zap.Logger, resp *response.Responder) *Middlewares {
	return &Middlewares{
		log:    log,
		resp:   resp,
		config: cfg,
	}
}
