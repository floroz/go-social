package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config      *config
	userService interfaces.UserService
}

type config struct {
	port string
}

func (app *application) run() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.port),
		Handler:      app.routes(),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Minute,
	}

	fmt.Printf("Starting server on %s\n", app.config.port)
	return server.ListenAndServe()
}

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/users", func(r chi.Router) {
			// this should be public
			r.Get("/", app.listUsersHandler)
			// all other routes should be protected
			r.Post("/", app.createUserHandler)
			r.Delete("/{id}", app.deleteUserHandler)
			r.Put("/{id}", app.updateUserHandler)
		})
	})

	return r
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	createUserDto := &domain.CreateUserDTO{}

	err := json.NewDecoder(r.Body).Decode(createUserDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := app.userService.Create(r.Context(), createUserDto)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)

}

func (app *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = parsedLimit
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil {
			offset = parsedOffset
		}
	}

	users, err := app.userService.List(r.Context(), limit, offset)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = app.userService.Delete(r.Context(), id)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	updateUserDto := &domain.UpdateUserDTO{}

	err = json.NewDecoder(r.Body).Decode(updateUserDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateUserDto.ID = id

	user, err := app.userService.Update(r.Context(), updateUserDto)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
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
