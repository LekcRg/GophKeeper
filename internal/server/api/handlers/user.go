package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	"go.uber.org/zap"
)

type UserService interface {
	Register(ctx context.Context, req models.RegisterUserReq) (models.TokenUserRes, error)
	Login(ctx context.Context, req models.LoginUserReq) (models.TokenUserRes, error)
}

type UserHandlers struct {
	resp    *response.Responder
	service UserService
	config  *config.Config
	log     *zap.Logger
}

func NewUserHandlers(
	cfg *config.Config, service UserService, log *zap.Logger, resp *response.Responder,
) *UserHandlers {
	return &UserHandlers{
		resp:    resp,
		service: service,
		config:  cfg,
		log:     log,
	}
}

func (uh *UserHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var reqBody models.RegisterUserReq

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		uh.resp.InternalError(w)

		return
	}
	defer r.Body.Close()

	res, err := uh.service.Register(r.Context(), reqBody)
	if err != nil {
		if errors.Is(err, errs.ErrLoginAlreadyExists) {
			uh.resp.JSON(w, http.StatusConflict, models.UserError{
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

func (uh *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var reqBody models.RegisterUserReq

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		uh.resp.InternalError(w)

		return
	}
	defer r.Body.Close()

	res, err := uh.service.Login(r.Context(), models.LoginUserReq(reqBody))
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			uh.resp.JSON(w, http.StatusMethodNotAllowed, models.UserError{
				Login: "invalid login or password",
			})

			return
		}

		uh.log.Error("UserService error", zap.Error(err))
		uh.resp.InternalError(w)

		return
	}

	uh.resp.JSON(w, http.StatusOK, res)
}
