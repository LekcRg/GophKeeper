package api

import (
	"github.com/LekcRg/GophKeeper/internal/server/api/handlers"
	"github.com/LekcRg/GophKeeper/internal/server/api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(h *handlers.Handlers, m *middlewares.Middlewares) *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		m.RequestLogger,
		middleware.CleanPath,
		middleware.AllowContentType("application/json"),
	)

	r.Route("/user", func(cr chi.Router) {
		cr.Post("/create", h.UserHandlers.Register)
		cr.Post("/login", h.UserHandlers.Login)
	})

	return r
}
