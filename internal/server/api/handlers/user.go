package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/api/middlewares"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
)

type UserService interface {
	Register(ctx context.Context, req models.UserReq) (models.APIKeyRes, error)
	GetCryptoParams(ctx context.Context, id int) (models.CryptoParamsRes, error)
	UpdateAPIKey(ctx context.Context, req models.UserLogin) (models.APIKeyRes, error)
	ChangePassword(ctx context.Context, req models.UserChangePasswordReq) error
}

type UserHandlers struct {
	resp    *response.Responder
	service UserService
	config  *config.Config
	log     *zap.Logger
}

type userReqService func(ctx context.Context, req models.UserReq) (any, error)

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

func userReqHandler[T any](
	w http.ResponseWriter,
	r *http.Request,
	log *zap.Logger,
	resp *response.Responder,
	handleServiceError func(http.ResponseWriter, error),
	serviceFunc func(context.Context, T) (models.APIKeyRes, error),
) {
	var reqBody T

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Error("Json decode error", zap.Error(err))
		resp.Error(w, http.StatusBadRequest, "Invalid JSON")

		return
	}
	defer r.Body.Close()

	res, err := serviceFunc(r.Context(), reqBody)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	resp.JSON(w, http.StatusCreated, res)
}

func (uh *UserHandlers) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrLoginAlreadyExists):
		uh.resp.JSON(w, http.StatusConflict, models.UserReq{
			Login: "Login already exists",
		})
	case errors.Is(err, errs.ErrInvalidCredentials):
		uh.resp.JSON(w, http.StatusBadRequest, models.UserReq{
			Login: "Invalid login or password",
		})
	case errors.Is(err, errs.ErrInvalidPassword):
		uh.resp.JSON(w, http.StatusBadRequest, models.UserChangePasswordReq{
			CurrentPassword: "Password is not correct",
		})
	case errors.As(err, &validation.Errors{}):
		uh.resp.JSON(w, http.StatusBadRequest, err)
	default:
		uh.log.Error("UserService unexpected error", zap.Error(err))
		uh.resp.InternalError(w)
	}
}

// Register godoc
// @Summary      Register user
// @Description  Register user and return API Key
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body models.UserReq true "Login and password"
// @Success      200 {object} models.APIKeyRes
// @Failure      400 {object} models.UserReq
// @Failure      409 {object} models.UserReq
// @Failure      500 {object} models.ResponseError
// @Router       /user/create [post]
//
// Register handles user registration and creating API Key.
func (uh *UserHandlers) Register(w http.ResponseWriter, r *http.Request) {
	userReqHandler(
		w, r,
		uh.log,
		uh.resp,
		uh.handleServiceError,
		uh.service.Register,
	)
}

// APIKey godoc
// @Summary      API Key generation
// @Description  Create or update user API Key
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body models.UserLogin true "Login and password"
// @Success      200 {object} models.APIKeyRes
// @Failure      400 {object} models.UserLogin
// @Failure      500 {object} models.ResponseError
// @Router       /user/api-key [post]
//
// APIKey handles create or update user API Key.
func (uh *UserHandlers) APIKey(w http.ResponseWriter, r *http.Request) {
	userReqHandler(
		w, r,
		uh.log,
		uh.resp,
		uh.handleServiceError,
		uh.service.UpdateAPIKey,
	)
}

// ChangePassword godoc
// @Summary      Change user password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body models.UserChangePasswordReq true "Current password and new password"
// @Success      200 {object} models.Response
// @Failure      400 {object} models.UserChangePasswordReq
// @Failure      500 {object} models.ResponseError
// @Router       /user/change-password [post]
// @Security     BearerAuth
//
// ChangePassword updates the authenticated user's password.
func (uh *UserHandlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req models.UserChangePasswordReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.log.Error("Json decode error", zap.Error(err))
		uh.resp.Error(w, http.StatusBadRequest, "Invalid JSON")

		return
	}

	err = uh.service.ChangePassword(r.Context(), req)
	if err != nil {
		uh.handleServiceError(w, err)

		return
	}

	uh.resp.JSON(w, http.StatusOK, models.Response{
		Message: "Password successfully changed",
	})
}

// GetCryptoParams godoc
// @Summary      Get crypto params
// @Description  Get crypto params
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200 {object} models.CryptoParamsRes
// @Failure      400 {object} models.UserLogin
// @Failure      409 {object} models.UserLogin
// @Failure      500 {object} models.ResponseError
// @Router       /user/crypto-params [get]
// @Security     BearerAuth
//
// GetCryptoParams handles user registration and creating API Key.
func (uh *UserHandlers) GetCryptoParams(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetID(r.Context())
	if err != nil {
		uh.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	res, err := uh.service.GetCryptoParams(r.Context(), id)
	if err != nil {
		uh.handleServiceError(w, err)

		return
	}

	uh.resp.JSON(w, http.StatusOK, res)
}
