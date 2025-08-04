package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/LekcRg/GophKeeper/internal/buildinfo"
	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/logger"
	"github.com/LekcRg/GophKeeper/internal/server/api"
	"github.com/LekcRg/GophKeeper/internal/server/api/handlers"
	"github.com/LekcRg/GophKeeper/internal/server/api/middlewares"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	"github.com/LekcRg/GophKeeper/internal/server/repository"
	"github.com/LekcRg/GophKeeper/internal/server/repository/postgres"
	"github.com/LekcRg/GophKeeper/internal/server/service"
	"github.com/LekcRg/GophKeeper/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	Log     *zap.Logger
	Config  *config.Config
	http    *http.Server
	db      *repository.Repository
	storage storage.Storage
}

var ErrLoggerIsNil = errors.New("logger is nil")

func NewServerApp(ctx context.Context) (*Server, error) {
	cfg, err := config.GetConfig(os.Args[1:])
	if err != nil || cfg == nil {
		return nil, err
	}

	log, err := logger.CreateServerLogger(cfg)
	if err != nil {
		return nil, err
	} else if log == nil {
		return nil, ErrLoggerIsNil
	}

	buildinfo.Print(log)

	server := &Server{
		Log:    log,
		Config: cfg,
	}

	server.printConfig()

	db, err := postgres.New(ctx, &cfg.Postgres, log)
	if err != nil {
		return nil, err
	}

	server.db = db

	server.storage, err = storage.New(cfg.Storage)
	if err != nil {
		return nil, err
	}

	server.http = server.createHTTP()

	return server, nil
}

func (s *Server) printConfig() {
	const redacted = "[REDACTED]"

	cfg := *s.Config
	cfg.Postgres.Password = redacted
	cfg.Storage.User = redacted
	cfg.Storage.Password = redacted

	s.Log.Info("Got config", zap.Any("config", cfg))
}

func (s *Server) createRouter() *chi.Mux {
	svc := service.New(s.db, s.Config, s.storage)
	resp := response.NewResponder(s.Log)
	handl := handlers.New(s.Config, svc, s.Log, resp)
	middl := middlewares.New(s.Config, s.Log, resp, s.db.UserRepo)

	return api.New(handl, middl)
}

func (s *Server) createHTTP() *http.Server {
	const (
		readTimeout       = 5 * time.Second
		writeTimeout      = 10 * time.Second
		readHeaderTimeout = 5 * time.Second
		idleTimeout       = 60 * time.Second
	)

	return &http.Server{
		Addr:              s.Config.Addr,
		Handler:           s.createRouter(),
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
	}
}

func (s *Server) startHTTPServer() error {
	s.Log.Info("Starting HTTP server", zap.String("HTTP address", s.http.Addr))

	err := s.http.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Start() error {
	return s.startHTTPServer()
}

func (s *Server) Shutdown(ctx context.Context) {
	s.Log.Info("Shutting down HTTP server...")

	err := s.http.Shutdown(ctx)
	if err != nil {
		s.Log.Warn("HTTP server shutdown error", zap.Error(err))

		if closeErr := s.http.Close(); closeErr != nil {
			s.Log.Error("HTTP server close error", zap.Error(closeErr))
		}
	} else {
		s.Log.Info("HTTP server gracefully stopped")
	}

	err = s.db.DB.Close()
	if err != nil {
		s.Log.Error("DB close error", zap.Error(err))
	} else {
		s.Log.Info("DB gracefully stopped")
	}
}
