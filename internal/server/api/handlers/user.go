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
	"go.uber.org/zap"
)

type UserService interface {
	Register(ctx context.Context, req models.UserReq) (models.TokenUserRes, error)
	Login(ctx context.Context, req models.UserReq) (models.TokenUserRes, error)
	ChangePassword(ctx context.Context, req models.UserChangePasswordReq) error
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

// Register godoc
// @Summary      Register user
// @Description  Register user and return JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body models.UserReq true "Login and password"
// @Success      200 {object} models.TokenUserRes
// @Failure      409 {object} models.UserError
// @Failure      500 {object} models.ResponseError
// @Router       /user/create [post]
//
// Register handles user registration.
func (uh *UserHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var reqBody models.UserReq

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

// Login godoc
// @Summary      Authentication
// @Description  Authentication and return JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body models.UserReq true "Login and password"
// @Success      200 {object} models.TokenUserRes
// @Failure      409 {object} models.UserError
// @Failure      500 {object} models.ResponseError
// @Router       /user/login [post]
//
// Login handles user authentication.
func (uh *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var reqBody models.UserReq

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		uh.resp.InternalError(w)

		return
	}
	defer r.Body.Close()

	res, err := uh.service.Login(r.Context(), models.UserReq(reqBody))
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

// ChangePassword godoc
// @Summary      Change user password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body models.UserChangePasswordReq true "Current password and new password"
// @Success      200 {object} models.Response
// @Failure      400 {object} models.ResponseError
// @Failure      500 {object} models.ResponseError
// @Router       /user/change-password [post]
// @Security     BearerAuth
//
// ChangePassword updates the authenticated user's password.
func (uh *UserHandlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req models.UserChangePasswordReq

	login, err := middlewares.GetLogin(r.Context())
	if err != nil {
		uh.resp.Error(w, http.StatusUnauthorized, "Unauthorized")

		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.resp.InternalError(w)
	}

	req.Login = login

	err = uh.service.ChangePassword(r.Context(), req)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidPassword) {
			uh.resp.Error(w, http.StatusBadRequest, "Invalid password")

			return
		}

		uh.resp.InternalError(w)

		return
	}

	uh.resp.JSON(w, http.StatusOK, models.Response{
		Message: "Password successfully changed",
	})
}
