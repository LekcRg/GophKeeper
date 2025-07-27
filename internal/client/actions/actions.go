package actions

import (
	"github.com/LekcRg/GophKeeper/internal/client/req"
	"go.uber.org/zap"
)

type Actions struct {
	req *req.Request
	log *zap.Logger
}

func New(request *req.Request, log *zap.Logger) *Actions {
	return &Actions{
		req: request,
		log: log,
	}
}
