package api

import (
	"net/http"
	"strconv"

	"github.com/floroz/go-social/internal/domain"
	"github.com/go-chi/chi/v5"
)

func (app *Application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	createUserDto := &domain.CreateUserDTO{}

	if err := readJSON(r.Body, createUserDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := app.UserService.Create(r.Context(), createUserDto)
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)

}

func (app *Application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
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

	users, err := app.UserService.List(r.Context(), limit, offset)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (app *Application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = app.UserService.Delete(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func (app *Application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	updateUserDto := &domain.UpdateUserDTO{}

	err = readJSON(r.Body, updateUserDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateUserDto.ID = id

	user, err := app.UserService.Update(r.Context(), updateUserDto)
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
