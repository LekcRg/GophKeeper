package storage

import "context"

type Storage interface {
	GenUploadPresignedURL(ctx context.Context, userID int) (url, path string, err error)
	GenPresignedGetURL(ctx context.Context, filePath string) (string, error)
	IsContainsFile(ctx context.Context, path string) bool
}
