package api

import (
	_ "github.com/LekcRg/GophKeeper/docs" // for swagger documentation.
	"github.com/LekcRg/GophKeeper/internal/server/api/handlers"
	"github.com/LekcRg/GophKeeper/internal/server/api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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

	r.Group(func(cr chi.Router) {
		cr.Use(m.Authenticate)

		cr.Post("/user/change-password", h.UserHandlers.ChangePassword)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return r
}
