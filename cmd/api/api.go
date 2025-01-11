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

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", app.loginHandler)
			r.Post("/signup", app.signupHandler)
			r.Post("/logout", app.logoutHandler)
			r.Post("/refresh", app.refreshHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware)
			r.Put("/", app.updateUserHandler)
			r.Get("/", app.getUserProfileHandler)
		})

		r.Route("/posts", func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware)

			r.Post("/", app.createPostHandler)
			r.Delete("/{id}", app.deletePostHandler)
			r.Put("/{id}", app.updatePostHandler)
			r.Get("/{id}", app.getPostByIdHandler)
			r.Get("/", app.listPostsHandler)

			r.Route("/{postId}/comments", func(r chi.Router) {
				r.Post("/", app.createCommentHandler)
				r.Put("/{id}", app.updateCommentHandler)
				r.Delete("/{id}", app.deleteCommentHandler)
				r.Get("/{id}", app.getCommentByIdHandler)
				r.Get("/", app.listByPostIdHandler)
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
		return
	}

	response := map[string]any{
		"data": data,
	}

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
		writeJSONError(w, http.StatusInternalServerError, "failed to write response")
	}
}

func readJSON(source io.Reader, dest any) error {
	// TODO: this could be extract into a middleware to specific routes if we need more fine-tuned control
	maxBytes := int64(1 << 20) // 1MB
	// Prevent reading too much data from the request body
	// to avoid potential denial of service attacks
	source = io.LimitReader(source, maxBytes)
	decoder := json.NewDecoder(source)

	// Ensure that the request body does not contain unknown fields
	decoder.DisallowUnknownFields()

	return decoder.Decode(dest)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	errors := []domain.ErrorDetail{
		{Message: message},
	}
	errorResponse := errorResponse{
		Errors: errors,
	}
	writeJSONResponse(w, status, errorResponse)
}

func handleErrors(w http.ResponseWriter, err error) {
	log.Error().Err(err).Msg("error")
	switch e := err.(type) {
	case *domain.ValidationError:
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case *domain.InternalServerError:
		writeJSONError(w, e.StatusCode, err.Error())
	case *domain.BadRequestError:
		writeJSONError(w, e.StatusCode, err.Error())
	case *domain.NotFoundError:
		writeJSONError(w, e.StatusCode, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
	}
}

func getUserClaimFromContext(ctx context.Context) (*domain.UserClaims, bool) {
	claims, ok := ctx.Value(middlewares.ContextKeyUser).(*domain.UserClaims)
	return claims, ok
}
