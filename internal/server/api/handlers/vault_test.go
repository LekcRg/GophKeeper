package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/api/middlewares"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestCreateItem(t *testing.T) {
	type test struct {
		svcErr   error
		ctx      context.Context
		body     models.VaultCreateItemReq
		name     string
		bodyStr  string
		wantErrs []string
		svcRes   models.VaultItem
		wantCode int
		id       int
		mockSvc  bool
	}

	successRes := models.VaultItem{
		Name:          "success",
		EncryptedData: []byte("data"),
	}

	ctx := context.Background()

	encData := []byte("data")

	tests := []test{
		{
			name: "Success password",
			body: models.VaultCreateItemReq{
				Name:          "test password",
				Type:          "password",
				EncryptedData: base64.StdEncoding.EncodeToString(encData),
			},
			svcRes:   successRes,
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			id:       1,
			wantCode: http.StatusOK,
			mockSvc:  true,
		},
		{
			name:     "invalid json",
			bodyStr:  `{name: name"}`,
			ctx:      ctx,
			wantCode: http.StatusBadRequest,
			wantErrs: []string{"error"},
		},
		{
			name: "Service validation error",
			body: models.VaultCreateItemReq{
				Name:          "test password",
				Type:          "invalid",
				EncryptedData: base64.StdEncoding.EncodeToString(encData),
			},
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			id:       1,
			wantCode: http.StatusBadRequest,
			wantErrs: []string{"type"},
			mockSvc:  true,
			svcErr: validation.Errors{
				"type": errs.ErrVaultNotCorrectType,
			},
		},
		{
			name: "Service internal error",
			body: models.VaultCreateItemReq{
				Name:          "test password",
				Type:          "invalid",
				EncryptedData: base64.StdEncoding.EncodeToString(encData),
			},
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			id:       1,
			wantCode: http.StatusInternalServerError,
			wantErrs: []string{"error"},
			mockSvc:  true,
			svcErr:   errors.New("internal"),
		},
		{
			name: "Invalid encrypted data",
			body: models.VaultCreateItemReq{
				Name:          "test password",
				Type:          "invalid",
				EncryptedData: "invalid",
			},
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			wantCode: http.StatusBadRequest,
			wantErrs: []string{"error"},
		},
		{
			name:     "Without ID",
			ctx:      ctx,
			wantCode: http.StatusUnauthorized,
			wantErrs: []string{"error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			require.NoError(t, err)

			vh, svc := getVaultHandlers(t)

			if tt.mockSvc {
				encData, err = base64.StdEncoding.DecodeString(tt.body.EncryptedData)
				bodySvc := models.VaultItem{
					UserID:        tt.id,
					Name:          tt.body.Name,
					Type:          tt.body.Type,
					EncryptedData: encData,
				}

				require.NoError(t, err)

				svc.EXPECT().CreateItem(tt.ctx, bodySvc).Return(successRes, tt.svcErr)
			}

			var bodyReader io.Reader = bytes.NewReader(body)
			if tt.bodyStr != "" {
				bodyReader = strings.NewReader(tt.bodyStr)
			}

			res := serveHTTPWithCtx(tt.ctx, vh.CreateItem, serveHTTPOpts{
				body: bodyReader,
			})

			assert.Equal(t, tt.wantCode, res.Code)

			if res.Code > 299 {
				checkSliceErrs(t, tt.wantErrs, res.Body.Bytes())
			} else {
				successResBytes, err := json.Marshal(tt.svcRes)
				require.NoError(t, err)

				assert.Equal(t, successResBytes, res.Body.Bytes())
			}
		})
	}
}

func TestGetAllItems(t *testing.T) {
	type test struct {
		svcErr   error
		ctx      context.Context
		name     string
		wantErrs []string
		wantCode int
		id       int
		mockSvc  bool
	}

	ctx := context.Background()
	tests := []test{
		{
			name:     "Success",
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			id:       1,
			wantCode: http.StatusOK,
			mockSvc:  true,
		},
		{
			name:     "Without ID",
			ctx:      ctx,
			wantCode: http.StatusUnauthorized,
			wantErrs: []string{"error"},
		},
		{
			name:     "Internal server error",
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			id:       1,
			wantCode: http.StatusInternalServerError,
			wantErrs: []string{"error"},
			svcErr:   errors.New("internal"),
			mockSvc:  true,
		},
	}

	listResSvc := []models.VaultItem{
		{
			Name: "name",
			Type: "note",
		},
		{
			Name: "name",
			Type: "password",
		},
	}

	listResHandler, err := json.Marshal(listResSvc)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vh, svc := getVaultHandlers(t)

			if tt.mockSvc {
				svc.EXPECT().GetAllItems(tt.ctx, tt.id).Return(listResSvc, tt.svcErr)
			}

			res := serveHTTPWithCtx(tt.ctx, vh.GetAllItems, serveHTTPOpts{
				method: "GET",
			})

			assert.Equal(t, tt.wantCode, res.Code)

			if res.Code > 299 {
				checkSliceErrs(t, tt.wantErrs, res.Body.Bytes())
			} else {
				assert.Equal(t, listResHandler, res.Body.Bytes())
			}
		})
	}
}

func TestCreateBinary(t *testing.T) {
	type test struct {
		ctx      context.Context
		svcErr   error
		body     models.VaultBinaryItemUploadReq
		name     string
		bodyStr  string
		encData  []byte
		wantErrs []string
		id       int
		wantCode int
		mockSvc  bool
	}

	encryptedData := []byte("data")
	encryptedDataStr := base64.StdEncoding.EncodeToString(encryptedData)
	ctx := context.Background()

	tests := []test{
		{
			name: "Success",
			body: models.VaultBinaryItemUploadReq{
				Name:          "test",
				EncryptedData: encryptedDataStr,
			},
			encData:  encryptedData,
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			id:       1,
			wantCode: http.StatusOK,
			mockSvc:  true,
		},
		{
			name:     "Without ID",
			ctx:      ctx,
			wantCode: http.StatusUnauthorized,
			wantErrs: []string{"error"},
		},
		{
			name:     "Internal server error",
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			id:       1,
			wantCode: http.StatusInternalServerError,
			wantErrs: []string{"error"},
			svcErr:   errors.New("internal"),
			mockSvc:  true,
		},
		{
			name:     "invalid json",
			bodyStr:  `{name: name"}`,
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			wantCode: http.StatusBadRequest,
			wantErrs: []string{"error"},
		},
		{
			name: "Invalid encrypted data",
			body: models.VaultBinaryItemUploadReq{
				Name:          "test",
				EncryptedData: "invalid",
			},
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			wantCode: http.StatusBadRequest,
			wantErrs: []string{"error"},
		},
		{
			name: "Validation error",
			body: models.VaultBinaryItemUploadReq{
				Name:          "",
				EncryptedData: encryptedDataStr,
			},
			encData:  encryptedData,
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			id:       1,
			wantCode: http.StatusBadRequest,
			wantErrs: []string{"name"},
			svcErr:   validation.Errors{"name": errors.New("name is required")},
			mockSvc:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			require.NoError(t, err)

			vh, svc := getVaultHandlers(t)

			svcRes := models.VaultBinaryItemUploadRes{
				URL:  "url",
				Path: "path",
				Item: models.VaultItem{
					Name:   tt.body.Name,
					UserID: tt.id,
				},
				ItemID: 1,
			}

			if tt.mockSvc {
				encData, err := base64.StdEncoding.DecodeString(tt.body.EncryptedData)
				bodySvc := models.VaultItem{
					UserID:        tt.id,
					Name:          tt.body.Name,
					Type:          "binary",
					EncryptedData: encData,
				}

				svcRes.Item.EncryptedData = encData

				require.NoError(t, err)

				svc.EXPECT().CreateBinary(tt.ctx, bodySvc).
					Return(svcRes, tt.svcErr)
			}

			var bodyReader io.Reader = bytes.NewReader(body)
			if tt.bodyStr != "" {
				bodyReader = strings.NewReader(tt.bodyStr)
			}

			res := serveHTTPWithCtx(tt.ctx, vh.CreateBinary, serveHTTPOpts{
				body: bodyReader,
			})

			assert.Equal(t, tt.wantCode, res.Code)

			if res.Code > 299 {
				checkSliceErrs(t, tt.wantErrs, res.Body.Bytes())
			} else {
				resJSON, err := json.Marshal(svcRes)
				require.NoError(t, err)
				assert.Equal(t, resJSON, res.Body.Bytes())
			}
		})
	}
}

func TestConfirmBinaryUpload(t *testing.T) {
	type test struct {
		ctx      context.Context
		svcErr   error
		body     models.VaultConfirmBinaryUploadReq
		name     string
		bodyStr  string
		wantErrs []string
		id       int
		wantCode int
		mockSvc  bool
	}

	ctx := context.Background()

	tests := []test{
		{
			name: "Success",
			body: models.VaultConfirmBinaryUploadReq{
				Path:    "path",
				VaultID: 1,
			},
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			id:       1,
			wantCode: http.StatusOK,
			mockSvc:  true,
		},
		{
			name:     "Internal server error",
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			id:       1,
			wantCode: http.StatusInternalServerError,
			wantErrs: []string{"error"},
			svcErr:   errors.New("internal"),
			mockSvc:  true,
		},
		{
			name:     "invalid json",
			bodyStr:  `{name: name"}`,
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			wantCode: http.StatusBadRequest,
			wantErrs: []string{"error"},
		},
		{
			name: "Validation error",
			body: models.VaultConfirmBinaryUploadReq{
				Path:    "",
				VaultID: 1,
			},
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			id:       1,
			wantCode: http.StatusBadRequest,
			wantErrs: []string{"path"},
			svcErr:   validation.Errors{"path": errors.New("path is required")},
			mockSvc:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			require.NoError(t, err)

			vh, svc := getVaultHandlers(t)

			if tt.mockSvc {
				svc.EXPECT().ConfirmBinaryUpload(tt.ctx, tt.body).
					Return(tt.svcErr)
			}

			var bodyReader io.Reader = bytes.NewReader(body)
			if tt.bodyStr != "" {
				bodyReader = strings.NewReader(tt.bodyStr)
			}

			res := serveHTTPWithCtx(tt.ctx, vh.ConfirmBinaryUpload, serveHTTPOpts{
				body: bodyReader,
			})

			assert.Equal(t, tt.wantCode, res.Code)

			if res.Code > 299 {
				checkSliceErrs(t, tt.wantErrs, res.Body.Bytes())
			}
		})
	}
}

func TestGetBinaryFileURL(t *testing.T) {
	type test struct {
		svcErr   error
		ctx      context.Context
		name     string
		vaultID  string
		wantErrs []string
		wantCode int
		userID   int
		mockSvc  bool
	}

	ctx := context.Background()
	tests := []test{
		{
			name:     "Success",
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			userID:   1,
			vaultID:  "2",
			wantCode: http.StatusOK,
			mockSvc:  true,
		},
		{
			name:     "Without ID",
			ctx:      ctx,
			wantCode: http.StatusUnauthorized,
			wantErrs: []string{"error"},
		},
		{
			name:     "Invalid vault id",
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			wantCode: http.StatusBadRequest,
			vaultID:  "",
			wantErrs: []string{"error"},
		},
		{
			name:     "Internal server error",
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			userID:   1,
			vaultID:  "2",
			mockSvc:  true,
			wantCode: http.StatusInternalServerError,
			wantErrs: []string{"error"},
			svcErr:   errors.New("internal"),
		},
		{
			name:     "Invalid user binary",
			ctx:      middlewares.AddIDToCtx(ctx, 1),
			userID:   1,
			vaultID:  "2",
			mockSvc:  true,
			wantCode: http.StatusBadRequest,
			wantErrs: []string{"error"},
			svcErr:   errs.ErrInvalidUserBinary,
		},
	}

	const url = "url"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vh, svc := getVaultHandlers(t)

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", tt.vaultID)
			ctx := context.WithValue(tt.ctx, chi.RouteCtxKey, chiCtx)

			if tt.mockSvc {
				vaultID, err := strconv.Atoi(tt.vaultID)
				require.NoError(t, err)
				svc.EXPECT().GetBinaryFileURL(ctx, tt.userID, vaultID).Return(url, tt.svcErr)
			}

			res := serveHTTPWithCtx(ctx, vh.GetBinaryFileURL, serveHTTPOpts{
				method: "GET",
			})

			assert.Equal(t, tt.wantCode, res.Code)

			if res.Code > 299 {
				checkSliceErrs(t, tt.wantErrs, res.Body.Bytes())
			} else {
				hRes, err := json.Marshal(models.GetBinaryFileURLRes{URL: url})
				require.NoError(t, err)
				assert.Equal(t, hRes, res.Body.Bytes())
			}
		})
	}
}

func getVaultHandlers(t *testing.T) (*VaultHandlers, *MockVaultService) {
	t.Helper()

	svc := NewMockVaultService(t)
	log := zaptest.NewLogger(t)
	resp := response.NewResponder(log)
	vh := NewVaultHandlers(&config.Config{}, svc, log, resp)

	return vh, svc
}

func checkSliceErrs(t *testing.T, expected []string, body []byte) {
	t.Helper()

	var er map[string]string
	err := json.Unmarshal(body, &er)
	require.NoError(t, err)

	for _, errName := range expected {
		assert.NotEmpty(t, er[errName])
	}

	for key, val := range er {
		if val != "" && !slices.Contains(expected, key) {
			t.Errorf("error {%s: %s} not contains in expected errs slice %v", key, val, expected)
		}
	}
}
