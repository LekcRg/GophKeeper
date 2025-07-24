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
		cr.Post("/api-key", h.UserHandlers.APIKey)
		cr.Post("/change-password", h.UserHandlers.ChangePassword)
	})

	r.Group(func(chiR chi.Router) {
		chiR.Use(m.Authenticate)
		chiR.Route("/vault", func(cr chi.Router) {
			cr.Post("/create", h.VaultHandlers.CreateItem)
		})
		chiR.Get("/user/crypto-params", h.UserHandlers.GetCryptoParams)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return r
}
