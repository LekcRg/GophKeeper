package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/LekcRg/GophKeeper/internal/app"
	"go.uber.org/zap"
)

// @title GophKeeper API
// @version 1.0
// @description Gophkeeper password manager HTTP API

// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	var wg sync.WaitGroup

	const ctxTimeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)

	sapp, err := app.NewServerApp(ctx)
	if err != nil {
		log.Fatal("New server app error ", err)
	}

	defer func() {
		cancel()

		if err = sapp.Log.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			sapp.Log.Error("Log sync error", zap.Error(err))
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		exitSignals(sapp)
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := sapp.Start(); err != nil {
			sapp.Log.Error("Server app error", zap.Error(err))

			if p, err := os.FindProcess(os.Getpid()); err == nil {
				err := p.Signal(syscall.SIGTERM)
				if err != nil {
					sapp.Log.Error("Process signall error", zap.Error(err))
				}
			}
		}
	}()

	wg.Wait()
}

func exitSignals(s *app.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-sigChan
	s.Log.Info("Received shutdown signal", zap.String("signal", sig.String()))

	const shutdownTimeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	defer cancel()
	s.Shutdown(ctx)
}
