package main

import (
	"github.com/LekcRg/GophKeeper/internal/client/views"
	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/logger"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

func main() {
	zapLog, err := logger.CreateClientLogger()
	if err != nil {
		panic(err)
	}

	cfg, err := config.GetClientConfig()
	if err != nil {
		zapLog.Info("Get config error", zap.Error(err))
	}

	defer func() {
		err := zapLog.Sync()
		if err != nil {
			zapLog.Error("logger.sync error", zap.Error(err))
		}
	}()

	v := views.New(zapLog, cfg)

	_, err = tea.NewProgram(v, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
