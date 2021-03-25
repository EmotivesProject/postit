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

		r.Route("/user", func(r chi.Router) {
			r.Post("/", createUser)
			r.Route("/{username}", func(r chi.Router) {
				r.Get("/", fetchUser)
			})
		})

		r.With(verifyJTW()).Route("/post", func(r chi.Router) {
			r.Get("/", fetchPost)
			r.Post("/", createPost)

			r.Route("/{post_id}", func(r chi.Router) {
				r.Delete("/", deletePost)

				r.Route("/like", func(r chi.Router) {
					r.Post("/", createLike)
					r.Route("/{like_id}", func(r chi.Router) {
						r.Delete("/", deleteLike)
					})
				})

				r.Route("/comment", func(r chi.Router) {
					r.Post("/", createComment)
					r.Route("/{comment_id}", func(r chi.Router) {
						r.Delete("/", deleteComment)
					})
				})
			})
		})

	})

	return r
}
