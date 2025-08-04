package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type storage struct {
	client *minio.Client
	config config.Storage
}

func New(cfg config.Storage) (Storage, error) {
	cl, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: cfg.Secure,
	})
	if err != nil {
		return nil, err
	}

	s := &storage{
		config: cfg,
		client: cl,
	}

	return s, nil
}

func (s *storage) GenUploadPresignedUrl(
	ctx context.Context, userID int,
) (url, path string, err error) {
	filename := uuid.New()
	filepath := fmt.Sprintf("%d/%s.enc", userID, filename)

	presigned, err := s.client.PresignedPutObject(
		ctx,
		s.config.Bucket,
		filepath,
		time.Hour,
	)
	if err != nil {
		return "", "", err
	}

	return presigned.String(), filepath, nil
}

func (s *storage) GenPresignedGetUrl(
	ctx context.Context, filePath string,
) (string, error) {
	presigned, err := s.client.PresignedGetObject(
		ctx,
		s.config.Bucket,
		filePath,
		time.Hour,
		nil,
	)
	if err != nil {
		return "", err
	}

	return presigned.String(), nil
}

func (s *storage) IsContainsFile(ctx context.Context, path string) bool {
	_, err := s.client.StatObject(ctx, s.config.Bucket, path, minio.StatObjectOptions{})
	if err != nil {
		return false
	}

	return true
}
