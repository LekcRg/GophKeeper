package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"go.uber.org/zap"
)

type UserService interface {
	Register(ctx context.Context, req models.RegisterUserReq) (models.RegisterUserRes, error)
}

type UserHandlers struct {
	resp    *Responder
	service UserService
	config  *config.Config
	log     *zap.Logger
}

func NewUserHandlers(
	cfg *config.Config, service UserService, log *zap.Logger, resp *Responder,
) *UserHandlers {
	return &UserHandlers{
		resp:    resp,
		service: service,
		config:  cfg,
		log:     log,
	}
}

func (uh *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var reqBody models.RegisterUserReq

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}
	defer r.Body.Close()

	res, err := uh.service.Register(r.Context(), reqBody)
	if err != nil {
		if errors.Is(err, errs.ErrLoginAlreadyExists) {
			uh.resp.JSON(w, http.StatusConflict, models.RegisterUserError{
				Login: "login already exists",
			})

			return
		}

		uh.log.Error("UserService error", zap.Error(err))
		uh.resp.InternalError(w)

		return
	}

	uh.resp.JSON(w, http.StatusOK, res)
}
