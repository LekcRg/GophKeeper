package main

import (
	"github.com/LekcRg/GophKeeper/internal/client/views"
	"github.com/LekcRg/GophKeeper/internal/logger"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

func main() {
	zapLog, err := logger.CreateClientLogger()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := zapLog.Sync()
		if err != nil {
			zapLog.Error("logger.sync error", zap.Error(err))
		}
	}()

	v := views.New(zapLog)

	_, err = tea.NewProgram(v, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
