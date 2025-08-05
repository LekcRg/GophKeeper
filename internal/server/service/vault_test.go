package service

import (
	"context"
	"errors"
	"testing"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/mocks"
	"github.com/LekcRg/GophKeeper/internal/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateItem(t *testing.T) {
	t.Parallel()

	cfg, err := config.GetConfig([]string{})
	require.NoError(t, err)

	type test struct {
		repoErr   error
		asErr     any
		name      string
		req       models.VaultItem
		wantErr   bool
		doNotMock bool
	}

	tests := []test{
		{
			name: "Success",
			req: models.VaultItem{
				Name:          "test",
				Type:          "password",
				EncryptedData: []byte("data"),
			},
		},
		{
			name: "Invalid type",
			req: models.VaultItem{
				Name:          "test",
				Type:          "invalid",
				EncryptedData: []byte("data"),
			},
			wantErr:   true,
			asErr:     validation.Errors{},
			doNotMock: true,
		},
		{
			name: "Internal repo err",
			req: models.VaultItem{
				Name:          "test",
				Type:          "password",
				EncryptedData: []byte("data"),
			},
			wantErr: true,
			repoErr: errors.New("internal"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockVaultRepo(t)
			st := mocks.NewMockStorage(t)
			vs := NewVaultService(repo, cfg, st)

			if !tt.doNotMock {
				repo.EXPECT().
					CreateItem(mock.Anything, tt.req).
					Return(tt.req, tt.repoErr)
			}

			item, err := vs.CreateItem(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)

				if tt.asErr != nil {
					assert.ErrorAs(t, err, &tt.asErr)
				}

				return
			}

			require.NoError(t, err)

			assert.Equal(t, tt.req, item)
		})
	}
}

func TestGetAllItems(t *testing.T) {
	t.Parallel()

	cfg, err := config.GetConfig([]string{})
	require.NoError(t, err)

	type test struct {
		repoErr   error
		asErr     any
		name      string
		id        int
		wantErr   bool
		doNotMock bool
	}

	tests := []test{
		{
			name: "Success",
			id:   1,
		},
		{
			name:    "Internal repo err",
			id:      1,
			wantErr: true,
			repoErr: errors.New("internal"),
		},
	}

	oneEncData := []byte("one")
	twoEncData := []byte("two")

	listRes := []models.VaultItem{
		{
			Name:          "One",
			Type:          "password",
			EncryptedData: oneEncData,
		},
		{
			Name:          "Two",
			Type:          "note",
			EncryptedData: twoEncData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockVaultRepo(t)
			st := mocks.NewMockStorage(t)
			vs := NewVaultService(repo, cfg, st)

			if !tt.doNotMock {
				cList := make([]models.VaultItem, len(listRes))
				copy(cList, listRes)

				repo.EXPECT().
					GetAllItems(mock.Anything, tt.id).
					Return(cList, tt.repoErr)
			}

			res, err := vs.GetAllItems(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)

				if tt.asErr != nil {
					assert.ErrorAs(t, err, &tt.asErr)
				}

				return
			}

			require.NoError(t, err)

			for _, item := range res {
				assert.NotEmpty(t, item.EncryptedDataString)
			}
		})
	}
}

func TestCreateBinary(t *testing.T) {
	t.Parallel()

	cfg, err := config.GetConfig([]string{})
	require.NoError(t, err)

	const (
		validType = "binary"                     // корректный тип
		mockURL   = "https://example.com/upload" // фиктивные данные для presign
		mockPath  = "users/1/abcd1234"
	)

	type test struct {
		repoErr          error
		storageErr       error
		asErr            any
		name             string
		req              models.VaultItem
		wantErr          bool
		doNotMockRepo    bool
		doNotMockStorage bool
	}

	tests := []test{
		{
			name: "Success",
			req: models.VaultItem{
				UserID:        1,
				Name:          "pic.jpg",
				Type:          validType,
				EncryptedData: []byte("data"),
			},
		},
		{
			name: "Invalid type",
			req: models.VaultItem{
				UserID:        1,
				Name:          "bad",
				Type:          "invalid",
				EncryptedData: []byte("data"),
			},
			wantErr:          true,
			asErr:            validation.Errors{},
			doNotMockRepo:    true, // валидация упадёт раньше
			doNotMockStorage: true, // до хранилища не дойдём
		},
		{
			name: "Internal repo err",
			req: models.VaultItem{
				UserID:        1,
				Name:          "doc.pdf",
				Type:          validType,
				EncryptedData: []byte("data"),
			},
			repoErr:          errors.New("internal"),
			wantErr:          true,
			doNotMockStorage: true, // хранилище не вызывается, т.к. repo вернул ошибку
		},
		{
			name: "Storage err",
			req: models.VaultItem{
				UserID:        2,
				Name:          "music.mp3",
				Type:          validType,
				EncryptedData: []byte("data"),
			},
			storageErr: errors.New("storage"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockVaultRepo(t)
			st := mocks.NewMockStorage(t)
			vs := NewVaultService(repo, cfg, st)

			createdItem := tt.req
			createdItem.ID = 42

			if !tt.doNotMockRepo {
				repo.EXPECT().
					CreateItem(mock.Anything, tt.req).
					Return(createdItem, tt.repoErr)
			}

			if !tt.doNotMockStorage {
				st.EXPECT().
					GenUploadPresignedURL(mock.Anything, tt.req.UserID).
					Return(mockURL, mockPath, tt.storageErr)
			}

			res, err := vs.CreateBinary(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)

				if tt.asErr != nil {
					assert.ErrorAs(t, err, &tt.asErr)
				}

				return
			}

			require.NoError(t, err)

			assert.Equal(t, createdItem, res.Item)
			assert.Equal(t, createdItem.ID, res.ItemID)
			assert.Equal(t, mockURL, res.URL)
			assert.Equal(t, mockPath, res.Path)
		})
	}
}

func TestConfirmBinaryUpload(t *testing.T) {
	t.Parallel()

	cfg, err := config.GetConfig([]string{})
	require.NoError(t, err)

	const (
		examplePath = "users/1/abcd1234"
	)

	type test struct {
		repoErr          error
		asErr            any
		name             string
		req              models.VaultConfirmBinaryUploadReq
		storageContains  bool
		wantErr          bool
		doNotMockRepo    bool
		doNotMockStorage bool
	}

	tests := []test{
		{
			name: "Success",
			req: models.VaultConfirmBinaryUploadReq{
				VaultID: 42,
				Path:    examplePath,
			},
			storageContains: true,
		},
		{
			name: "Validation error (empty path)",
			req: models.VaultConfirmBinaryUploadReq{
				VaultID: 42,
				Path:    "",
			},
			wantErr:          true,
			asErr:            validation.Errors{},
			doNotMockRepo:    true,
			doNotMockStorage: true,
		},
		{
			name: "Binary not found in storage",
			req: models.VaultConfirmBinaryUploadReq{
				VaultID: 42,
				Path:    examplePath,
			},
			storageContains: false,
			wantErr:         true,
			asErr:           errs.ErrBinaryFileNotFound,
			doNotMockRepo:   true,
		},
		{
			name: "Repo update error",
			req: models.VaultConfirmBinaryUploadReq{
				VaultID: 42,
				Path:    examplePath,
			},
			storageContains: true,
			repoErr:         errors.New("internal"),
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockVaultRepo(t)
			st := mocks.NewMockStorage(t)
			vs := NewVaultService(repo, cfg, st)

			if !tt.doNotMockStorage {
				st.EXPECT().
					IsContainsFile(mock.Anything, tt.req.Path).
					Return(tt.storageContains)
			}

			if !tt.doNotMockRepo {
				repo.EXPECT().
					UpdateBinaryURL(mock.Anything, tt.req).
					Return(tt.repoErr)
			}

			err := vs.ConfirmBinaryUpload(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)

				if tt.asErr != nil {
					assert.ErrorAs(t, err, &tt.asErr)
				}

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGetBinaryFileURL(t *testing.T) {
	t.Parallel()

	cfg, err := config.GetConfig([]string{})
	require.NoError(t, err)

	const (
		binaryPath = "users/1/file.bin"
		presignURL = "https://example.com/get"
	)

	type test struct {
		repoErr          error
		storageErr       error
		asErr            any
		name             string
		item             models.VaultItem
		userID           int
		vaultID          int
		wantErr          bool
		doNotMockRepo    bool
		doNotMockStorage bool
	}

	tests := []test{
		{
			name:    "Success",
			userID:  1,
			vaultID: 10,
			item: models.VaultItem{
				ID:         10,
				UserID:     1,
				BinaryPath: binaryPath,
			},
		},
		{
			name:             "Repo error",
			userID:           1,
			vaultID:          10,
			repoErr:          errors.New("db"),
			wantErr:          true,
			doNotMockStorage: true,
		},
		{
			name:    "Invalid user (not owner)",
			userID:  2,
			vaultID: 10,
			item: models.VaultItem{
				ID:         10,
				UserID:     1,
				BinaryPath: binaryPath,
			},
			wantErr:          true,
			asErr:            errs.ErrInvalidUserBinary,
			doNotMockStorage: true,
		},
		{
			name:    "Storage presign error",
			userID:  1,
			vaultID: 10,
			item: models.VaultItem{
				ID:         10,
				UserID:     1,
				BinaryPath: binaryPath,
			},
			storageErr: errors.New("s3 fail"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockVaultRepo(t)
			st := mocks.NewMockStorage(t)
			vs := NewVaultService(repo, cfg, st)

			if !tt.doNotMockRepo {
				repo.EXPECT().
					GetItem(mock.Anything, tt.vaultID).
					Return(tt.item, tt.repoErr)
			}

			if !tt.doNotMockStorage {
				st.EXPECT().
					GenPresignedGetURL(mock.Anything, tt.item.BinaryPath).
					Return(presignURL, tt.storageErr)
			}

			url, err := vs.GetBinaryFileURL(context.Background(), tt.userID, tt.vaultID)
			if tt.wantErr {
				assert.Error(t, err)

				if tt.asErr != nil {
					assert.ErrorAs(t, err, &tt.asErr)
				}

				return
			}

			require.NoError(t, err)
			assert.Equal(t, presignURL, url)
		})
	}
}
