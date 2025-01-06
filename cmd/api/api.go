package api

import (
	"net/http"

	"github.com/floroz/go-social/internal/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	Config      *Config
	UserService interfaces.UserService
	PostService interfaces.PostService
}

type Config struct {
	Port string
}

func (app *Application) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/users", func(r chi.Router) {
			// public
			r.Post("/", app.createUserHandler)
			// protected
			r.Delete("/{id}", app.deleteUserHandler)
			r.Put("/{id}", app.updateUserHandler)
			r.Get("/", app.listUsersHandler)
		})

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Delete("/{id}", app.deletePostHandler)
			r.Put("/{id}", app.updatePostHandler)
			r.Get("/{id}", app.getPostByIdHandler)
			r.Get("/", app.listPostsHandler)
		})
	})

	return r
}
