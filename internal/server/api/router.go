package api

import (
	"github.com/LekcRg/GophKeeper/internal/server/api/handlers"
	"github.com/go-chi/chi/v5"
)

func New(h *handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/user", func(cr chi.Router) {
		cr.Get("/create", h.UserHandlers.CreateUser)
	})

	return r
}
