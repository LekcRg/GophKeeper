package middlewares

import (
	"go.uber.org/zap"
)

type Middlewares struct {
	log *zap.Logger
}

func New(log *zap.Logger) *Middlewares {
	return &Middlewares{
		log: log,
	}
}
