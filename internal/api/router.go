package api

import (
	"github.com/go-chi/chi"
)

// can also do r.Use(authorizationMiddleware()) and r.With(authorization)
func CreateRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(SimpleMiddleware())

	r.Route("/postit/", func(r chi.Router) {
		r.Get("/healthz", healthz)

	})

	return r
}
