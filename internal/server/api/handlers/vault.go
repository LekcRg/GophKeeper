package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/api/middlewares"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
)

type VaultService interface {
	CreateItem(ctx context.Context, item models.VaultItem) (models.VaultItem, error)
	GetAllItems(ctx context.Context, id int) ([]models.VaultItem, error)
}

type VaultHandlers struct {
	resp    *response.Responder
	service VaultService
	config  *config.Config
	log     *zap.Logger
}

func NewVaultHandlers(
	cfg *config.Config, service VaultService, log *zap.Logger, resp *response.Responder,
) *VaultHandlers {
	return &VaultHandlers{
		resp:    resp,
		service: service,
		config:  cfg,
		log:     log,
	}
}

func (vh *VaultHandlers) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req models.VaultCreateItemReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		vh.resp.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	defer r.Body.Close()

	userID, err := middlewares.GetID(r.Context())
	if err != nil {
		vh.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	enctyptedBytes, err := base64.StdEncoding.DecodeString(req.EncryptedData)
	if err != nil {
		vh.resp.Error(w, http.StatusBadRequest, "Invalid encrypted data")
		return
	}

	vaultItem, err := vh.service.CreateItem(r.Context(), models.VaultItem{
		UserID:        userID,
		Name:          req.Name,
		Type:          req.Type,
		EncryptedData: enctyptedBytes,
	})
	if err != nil {
		if errors.As(err, &validation.Errors{}) {
			vh.resp.JSON(w, http.StatusBadRequest, err)

			return
		}

		vh.log.Info("Create vault service error", zap.Error(err))
		vh.resp.InternalError(w)

		return
	}

	vh.resp.JSON(w, http.StatusOK, vaultItem)
}

func (vh *VaultHandlers) GetAllItems(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetID(r.Context())
	if err != nil {
		vh.resp.Error(w, http.StatusBadRequest, "Unauthorized")
		return
	}

	res, err := vh.service.GetAllItems(r.Context(), id)
	if err != nil {
		vh.resp.InternalError(w)
		return
	}

	vh.resp.JSON(w, http.StatusOK, res)
}
