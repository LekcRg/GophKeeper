package logger

import (
	"path/filepath"

	"go.uber.org/zap"
)

func CreateClientLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	path, err := filepath.Abs("./output.log")
	if err != nil {
		return nil, err
	}

	cfg.OutputPaths = []string{
		path,
	}

	return cfg.Build()
}
