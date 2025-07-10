package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/models"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(ctx context.Context) error
	Test(ctx context.Context, message string) error
}

type UserHandlers struct {
	service UserService
	config  *config.Config
	log     *zap.Logger
}

func NewUserHandlers(cfg *config.Config, service UserService, log *zap.Logger) *UserHandlers {
	return &UserHandlers{
		service: service,
		config:  cfg,
		log:     log,
	}
}

func (uh *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var reqBody models.CreateUserReq

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	uh.log.Info("Body", zap.Any("CreateUser", reqBody))
	uh.log.Info("", zap.String("Content-Type", r.Header.Get("Content-Type")))

	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)

	_, err = io.WriteString(w, "Hello from user handler")
	if err != nil {
		uh.log.Error("CreateUser error", zap.Error(err))
	}
}

func (uh *UserHandlers) Test(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Test string `json:"test"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		uh.log.Error("Json decode error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	uh.log.Info(body.Test)

	defer r.Body.Close()

	res, err := json.Marshal(map[string]string{
		"status": "ok",
	})
	if err != nil {
		uh.log.Error("Json marshal error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	err = uh.service.Test(r.Context(), body.Test)
	if err != nil {
		uh.log.Error("UserService error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
