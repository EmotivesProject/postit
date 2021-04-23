package api

import (
	"github.com/TomBowyerResearchProject/common/middlewares"
	"github.com/TomBowyerResearchProject/common/verification"
	"github.com/go-chi/chi"
)

func CreateRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.SimpleMiddleware())

	r.Route("/", func(r chi.Router) {
		r.Get("/healthz", healthz)

		r.Route("/user", func(r chi.Router) {
			r.With(verification.VerifyJTW()).Get("/", fetchUserFromAuth)
		})

		r.With(verification.VerifyJTW()).Route("/post", func(r chi.Router) {
			r.Post("/", createPost)
			r.Get("/", fetchPosts)

			r.Route("/{post_id}", func(r chi.Router) {
				r.Get("/", fetchIndividualPost)

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
