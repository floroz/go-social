package api

import (
	"encoding/json"
	"net/http"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	Config      *Config
	UserService interfaces.UserService
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

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/users", func(r chi.Router) {
			// PUBLIC
			r.Post("/", app.createUserHandler)
			// PRIVATE
			r.Delete("/{id}", app.deleteUserHandler)
			r.Put("/{id}", app.updateUserHandler)
			r.Get("/", app.listUsersHandler)
		})
	})

	return r
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	data := map[string]string{"error": message}
	writeJSON(w, status, data)
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *domain.ValidationError:
		errorResponse(w, http.StatusBadRequest, e.Error())
	case *domain.ConflictError:
		errorResponse(w, http.StatusConflict, e.Error())
	case *domain.InternalServerError:
		errorResponse(w, http.StatusInternalServerError, e.Error())
	case *domain.BadRequestError:
		errorResponse(w, http.StatusBadRequest, e.Error())
	default:
		errorResponse(w, http.StatusInternalServerError, "internal server error")
	}
}
