package api

import (
	"github.com/go-chi/chi"
)

// can also do r.Use(authorizationMiddleware()) and r.With(authorization)
func CreateRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(SimpleMiddleware())

	r.Route("/", func(r chi.Router) {
		r.Get("/healthz", healthz)

		r.Route("/post", func(r chi.Router) {
			r.Get("/", fetchPost)
			r.Post("/", createPost)

			r.Route("/{post_id}", func(r chi.Router) {
				r.Delete("/", deletePost)
			})
		})

	})

	return r
}
