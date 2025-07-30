package api

import (
	_ "github.com/LekcRg/GophKeeper/docs" // for swagger documentation.
	"github.com/LekcRg/GophKeeper/internal/routes"
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

	r.Post(routes.UserRegister, h.UserHandlers.Register)
	r.Post(routes.UserUpdateKey, h.UserHandlers.APIKey)
	r.Post(routes.UserChangePassword, h.UserHandlers.ChangePassword)

	r.Group(func(cr chi.Router) {
		cr.Use(m.Authenticate)
		cr.Get(routes.UserGetCryptoParams, h.UserHandlers.GetCryptoParams)
		cr.Post(routes.VaultCreateItem, h.VaultHandlers.CreateItem)
		cr.Get(routes.VaultGetAll, h.VaultHandlers.GetAllItems)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return r
}
