package storage

import "context"

type Storage interface {
	GenUploadPresignedUrl(ctx context.Context, userID int) (url, path string, err error)
	GenPresignedGetUrl(ctx context.Context, filePath string) (string, error)
	IsContainsFile(ctx context.Context, path string) bool
}
