package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/floroz/go-social/cmd/middlewares"
	"github.com/floroz/go-social/internal/apitypes"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/errorcodes"
	"github.com/floroz/go-social/internal/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
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

		// Serve Swagger UI and OpenAPI spec
		apiRouter.Get("/docs", app.serveDocsHandler)
		apiRouter.Get("/openapi.yaml", app.serveOpenapiHandler)

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
		// Pass "" for fieldName as it's an internal server error writing response
		writeJSONError(w, http.StatusInternalServerError, "failed to write response", errorcodes.CodeInternalServerError, "")
	}
}

func readJSON(source io.Reader, dest any) error {
	maxBytes := int64(1 << 20) // 1MB
	source = io.LimitReader(source, maxBytes)
	decoder := json.NewDecoder(source)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dest)
}

// Modified to include fieldName
func writeJSONError(w http.ResponseWriter, status int, message string, code errorcodes.ApiErrorCode, fieldName string) {
	apiErr := apitypes.ApiError{
		Code:    string(code),
		Message: message,
	}
	if fieldName != "" {
		apiErr.Field = &fieldName
	}

	errorResponse := apitypes.ApiErrorResponse{
		Errors: []apitypes.ApiError{apiErr},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Error().Err(err).Msg("failed to write error response")
	}
}

func handleErrors(w http.ResponseWriter, err error) {
	log.Error().Err(err).Msg("handling error")

	// Check for specific sentinel errors first
	if errors.Is(err, domain.ErrNotFound) {
		writeJSONError(w, http.StatusNotFound, err.Error(), errorcodes.CodeNotFound, "")
		return
	}
	if errors.Is(err, domain.ErrDuplicateEmailOrUsername) { // New check for conflict
		writeJSONError(w, http.StatusConflict, err.Error(), errorcodes.CodeConflict, "")
		return
	}

	// Then check for custom error types
	switch e := err.(type) {
	case validator.ValidationErrors:
		if len(e) > 0 {
			// For now, take the first validation error to populate the single ApiError.
			// A more robust solution might involve returning multiple ApiError entries.
			firstErr := e[0]
			// Convert struct field name (PascalCase) to snake_case for JSON field name
			fieldName := toSnakeCase(firstErr.Field())
			// The message from firstErr.Error() is usually quite informative.
			writeJSONError(w, http.StatusBadRequest, firstErr.Error(), errorcodes.CodeValidationError, fieldName)
		} else {
			// Should not happen if validator.ValidationErrors is not empty, but as a fallback:
			writeJSONError(w, http.StatusBadRequest, "Validation failed", errorcodes.CodeValidationError, "")
		}
	case *domain.ValidationError: // This case might become less common if services return validator.ValidationErrors directly
		// If domain.ValidationError.Field is set meaningfully, pass it. Otherwise, pass ""
		field := ""
		if e.Field != "" && e.Field != "validation" { // Avoid passing "validation" as the field name
			field = e.Field
		}
		writeJSONError(w, http.StatusBadRequest, e.Error(), errorcodes.CodeValidationError, field)
	case *domain.InternalServerError:
		writeJSONError(w, e.StatusCode, e.Error(), errorcodes.CodeInternalServerError, "")
	case *domain.BadRequestError:
		writeJSONError(w, e.StatusCode, e.Error(), errorcodes.CodeBadRequest, "")
	case *domain.NotFoundError:
		writeJSONError(w, e.StatusCode, e.Error(), errorcodes.CodeNotFound, "")
	case *domain.ForbiddenError:
		writeJSONError(w, e.StatusCode, e.Error(), errorcodes.CodeForbidden, "")
	default:
		// Fallback for other unknown errors
		writeJSONError(w, http.StatusInternalServerError, "An unexpected internal server error occurred.", errorcodes.CodeInternalServerError, "")
	}
}

func getUserClaimFromContext(ctx context.Context) (*domain.UserClaims, bool) {
	claims, ok := ctx.Value(middlewares.ContextKeyUser).(*domain.UserClaims)
	return claims, ok
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// serveDocsHandler serves the Swagger UI HTML page.
func (app *Application) serveDocsHandler(w http.ResponseWriter, r *http.Request) {
	// Construct the absolute path to docs/index.html relative to the project root
	// Assuming the executable runs from the project root. Adjust if needed.
	cwd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current working directory")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(cwd, "docs", "index.html")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Error().Str("path", filePath).Msg("Swagger UI HTML file not found")
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}

// serveOpenapiHandler serves the bundled OpenAPI specification file.
func (app *Application) serveOpenapiHandler(w http.ResponseWriter, r *http.Request) {
	// Construct the absolute path to openapi/openapi-bundled.yaml
	cwd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current working directory")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(cwd, "openapi", "openapi-bundled.yaml")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Error().Str("path", filePath).Msg("OpenAPI spec file not found")
		http.NotFound(w, r)
		return
	}

	// Set appropriate content type
	w.Header().Set("Content-Type", "application/vnd.oai.openapi+yaml;charset=utf-8")
	http.ServeFile(w, r, filePath)
}
