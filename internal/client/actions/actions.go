package actions

import (
	"github.com/LekcRg/GophKeeper/internal/client/req"
	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/models"
	"go.uber.org/zap"
)

type Actions struct {
	req    *req.Request
	log    *zap.Logger
	config *config.ClientConfig
}

func New(request *req.Request, log *zap.Logger, cfg *config.ClientConfig) *Actions {
	return &Actions{
		req:    request,
		log:    log,
		config: cfg,
	}
}

func (a *Actions) UpdateConfigAddress(addr string) error {
	return a.config.Update(func(cfg *config.ClientConfig) {
		cfg.Address = addr
	})
}

func (a *Actions) UpdateConfigCredentials(c models.ClientRegisterResponse) error {
	err := a.config.Update(func(cfg *config.ClientConfig) {
		cfg.EnctyptedTag = c.Tag
		cfg.Salt = c.Salt
		cfg.Key = c.Key
	})

	return err
}

func (a *Actions) SaveConfig(config.ClientConfig) error {
	return nil
}
