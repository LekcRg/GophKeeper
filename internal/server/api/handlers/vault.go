package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/api/middlewares"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
)

type VaultService interface {
	CreateItem(ctx context.Context, item models.VaultItem) (models.VaultItem, error)
	GetAllItems(ctx context.Context, id int) ([]models.VaultItem, error)
	CreateBinary(ctx context.Context, req models.VaultItem) (models.VaultBinaryItemUploadRes, error)
	ConfirmBinaryUpload(ctx context.Context, req models.VaultConfirmBinaryUploadReq) error
	GetBinaryFileURL(ctx context.Context, userID, vaultID int) (string, error)
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

// CreateItem godoc
// @Summary      Create vault item
// @Description  Creates a new vault item (text, login, card, etc.) for the authenticated user.
// @Tags         Vault
// @Accept       json
// @Produce      json
// @Param        request body models.VaultCreateItemReq true "Vault item to create"
// @Success      201 {object} models.VaultItem
// @Failure      400 {object} models.VaultItem
// @Failure      401 {object} models.ResponseError
// @Failure      500 {object} models.ResponseError
// @Router       /vault [post]
// @Security     BearerAuth
//
// CreateItem handles creation of a new vault item (e.g., text, login, card) for the authenticated user.
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

// GetAllItems godoc
// @Summary      Get all vault items
// @Description  Returns all vault items belonging to the authenticated user.
// @Tags         Vault
// @Produce      json
// @Success      200 {array} models.VaultItem
// @Failure      401 {object} models.ResponseError
// @Failure      500 {object} models.ResponseError
// @Router       /vault/all [get]
// @Security     BearerAuth
//
// GetAllItems returns all vault items that belong to the authenticated user.
func (vh *VaultHandlers) GetAllItems(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetID(r.Context())
	if err != nil {
		vh.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	res, err := vh.service.GetAllItems(r.Context(), id)
	if err != nil {
		vh.log.Error("Vault get all error", zap.Error(err))
		vh.resp.InternalError(w)

		return
	}

	vh.resp.JSON(w, http.StatusOK, res)
}

// CreateBinary godoc
// @Summary      Create binary vault item
// @Description  Creates a new vault item of type 'binary' and returns a pre-signed URL for direct file upload.
// @Tags         Vault
// @Accept       json
// @Produce      json
// @Param        request body models.VaultBinaryItemUploadReq true "Binary vault item metadata"
// @Success      201 {object} models.VaultBinaryItemUploadRes
// @Failure      400 {object} models.VaultBinaryItemUploadRes
// @Failure      401 {object} models.ResponseError
// @Failure      500 {object} models.ResponseError
// @Router       /vault/binary [post]
// @Security     BearerAuth
//
// Creates a new vault item of type 'binary' and returns a pre-signed URL for direct file upload.
func (vh *VaultHandlers) CreateBinary(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetID(r.Context())
	if err != nil {
		vh.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.VaultBinaryItemUploadReq

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		vh.resp.Error(w, http.StatusBadRequest, "Bad request")
		return
	}

	defer r.Body.Close()

	enctyptedBytes, err := base64.StdEncoding.DecodeString(req.EncryptedData)
	if err != nil {
		vh.resp.Error(w, http.StatusBadRequest, "Invalid encrypted data")
		return
	}

	res, err := vh.service.CreateBinary(r.Context(), models.VaultItem{
		UserID:        id,
		Name:          req.Name,
		Type:          "binary",
		EncryptedData: enctyptedBytes,
	})
	if err != nil {
		if errors.As(err, &validation.Errors{}) {
			vh.resp.JSON(w, http.StatusBadRequest, err)

			return
		}

		vh.log.Error("Vault create binary error", zap.Error(err))
		vh.resp.InternalError(w)

		return
	}

	vh.resp.JSON(w, http.StatusOK, res)
}

// ConfirmBinaryUpload godoc
// @Summary      Confirm binary file upload
// @Description  Confirms the successful upload of a binary file and finalizes the vault item.
// @Tags         Vault
// @Accept       json
// @Produce      json
// @Param        request body models.VaultConfirmBinaryUploadReq true "Binary upload confirmation request"
// @Success      200 {object} models.ResponseMessage
// @Failure      400 {object} models.VaultConfirmBinaryUploadReq
// @Failure      401 {object} models.ResponseError
// @Failure      500 {object} models.ResponseError
// @Router       /vault/binary/confirm [post]
// @Security     BearerAuth
//
// Finalizes the binary vault item after successful file upload using the provided file path and item ID.
func (vh *VaultHandlers) ConfirmBinaryUpload(w http.ResponseWriter, r *http.Request) {
	var req models.VaultConfirmBinaryUploadReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		vh.resp.Error(w, http.StatusBadRequest, "Bad request")
		return
	}

	defer r.Body.Close()

	err = vh.service.ConfirmBinaryUpload(r.Context(), req)
	if err != nil {
		if errors.As(err, &validation.Errors{}) {
			vh.resp.JSON(w, http.StatusBadRequest, err)

			return
		}

		vh.log.Error("Vault confirm binary upload error",
			zap.Error(err), zap.Any("req", req))
		vh.resp.InternalError(w)

		return
	}

	vh.resp.Message(w, http.StatusOK, "Success")
}

// ConfirmBinaryUpload godoc
// @Summary      Get binary file URL
// @Description  Get binary file URL
// @Tags         Vault
// @Produce      json
// @Param        id   path      string  true  "Binary ID"
// @Success      200 {object} models.ResponseMessage
// @Failure      400 {object} models.ResponseError
// @Failure      401 {object} models.ResponseError
// @Failure      500 {object} models.ResponseError
// @Router       /vault/binary/{id} [get]
// @Security     BearerAuth
//
// Get binary file URL.
func (vh *VaultHandlers) GetBinaryFileURL(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.GetID(r.Context())
	if err != nil {
		vh.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vaultIDStr := chi.URLParam(r, "id")

	vaultID, err := strconv.Atoi(vaultIDStr)
	if err != nil {
		vh.resp.Error(w, http.StatusBadRequest, "Ivalid ID")
	}

	// vh.log.Info("got resp", zap.Int("user id", userID), zap.Int("vault id", vaultID))
	url, err := vh.service.GetBinaryFileURL(r.Context(), userID, vaultID)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidUserBinary) {
			vh.resp.Error(w, http.StatusBadRequest, "Bad request")
		}

		vh.log.Error("GetBinaryFileURL service error", zap.Error(err))
		vh.resp.InternalError(w)

		return
	}

	vh.resp.JSON(w, http.StatusOK, models.GetBinaryFileURLRes{URL: url})
}
