package logger

import (
	"github.com/LekcRg/GophKeeper/internal/config"
	"go.uber.org/zap"
)

func CreateLogger(cfg *config.Config) (*zap.Logger, error) {
	var (
		logger *zap.Logger
		err    error
	)

	if cfg.IsDev {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, err
	}

	return logger, nil
}
