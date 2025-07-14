package buildinfo

import (
	"go.uber.org/zap"
)

//nolint:gochecknoglobals // build info
var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func Print(log *zap.Logger) {
	log.Info("Build info",
		zap.String("version", BuildVersion),
		zap.String("data", BuildDate),
		zap.String("commit", BuildCommit),
	)
}
