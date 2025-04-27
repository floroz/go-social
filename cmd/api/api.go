package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/floroz/go-social/cmd/middlewares"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors" // Import cors package
	"github.com/rs/zerolog/log"
)

type Application struct {
	Config         *Config
	AuthService    interfaces.AuthService
	UserService    interfaces.UserService
	PostService    interfaces.PostService
	CommentService interfaces.CommentService
}

type Config struct {
	Port string
}

func (app *Application) Routes() http.Handler {
	r := chi.NewRouter()

	// Add CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"}, // Allow frontend dev server
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Get("/healthz", app.healthCheckHandler)

		// Versioned resource routes under /api/v1
		apiRouter.Route("/v1", func(v1Router chi.Router) {
			// Auth routes
			v1Router.Route("/auth", func(authRouter chi.Router) {
				authRouter.Post("/login", app.loginHandler)
				authRouter.Post("/signup", app.signupHandler)
				authRouter.Post("/logout", app.logoutHandler)
				authRouter.Post("/refresh", app.refreshHandler)
			})

			// User routes
			v1Router.Route("/users", func(userRouter chi.Router) {
				userRouter.Use(middlewares.AuthMiddleware)
				userRouter.Put("/", app.updateUserHandler)
				userRouter.Get("/", app.getUserProfileHandler)
			})

			// Post routes
			v1Router.Route("/posts", func(postRouter chi.Router) {
				postRouter.Use(middlewares.AuthMiddleware)
				postRouter.Post("/", app.createPostHandler)
				postRouter.Delete("/{id}", app.deletePostHandler)
				postRouter.Put("/{id}", app.updatePostHandler)
				postRouter.Get("/{id}", app.getPostByIdHandler)
				postRouter.Get("/", app.listPostsHandler)

				// Comments sub-route
				postRouter.Route("/{postId}/comments", func(commentRouter chi.Router) {
					// Auth middleware is already applied by the parent /posts route
					commentRouter.Post("/", app.createCommentHandler)
					commentRouter.Put("/{id}", app.updateCommentHandler)
					commentRouter.Delete("/{id}", app.deleteCommentHandler)
					commentRouter.Get("/{id}", app.getCommentByIdHandler)
					commentRouter.Get("/", app.listByPostIdHandler)
				})
			})
		})
	})

	return r
}

type errorResponse struct {
	Errors []domain.ErrorDetail `json:"errors"`
}

func writeJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		// If data is nil, we've already set the status code, just return.
		return
	}
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
		writeJSONError(w, http.StatusInternalServerError, "failed to write response")
	}
}

func readJSON(source io.Reader, dest any) error {
	maxBytes := int64(1 << 20) // 1MB
	source = io.LimitReader(source, maxBytes)
	decoder := json.NewDecoder(source)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dest)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	errors := []domain.ErrorDetail{{Message: message}}
	errorResponse := errorResponse{Errors: errors}
	// Use writeJSONResponse to ensure consistent { "data": { "errors": [...] } } structure for errors too
	writeJSONResponse(w, status, errorResponse)
}

func handleErrors(w http.ResponseWriter, err error) {
	log.Error().Err(err).Msg("error")
	switch e := err.(type) {
	case *domain.ValidationError:
		writeJSONError(w, http.StatusBadRequest, e.Error())
	case *domain.InternalServerError:
		writeJSONError(w, e.StatusCode, e.Error())
	case *domain.BadRequestError:
		writeJSONError(w, e.StatusCode, e.Error())
	case *domain.NotFoundError:
		writeJSONError(w, e.StatusCode, e.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
	}
}

func getUserClaimFromContext(ctx context.Context) (*domain.UserClaims, bool) {
	claims, ok := ctx.Value(middlewares.ContextKeyUser).(*domain.UserClaims)
	return claims, ok
}
