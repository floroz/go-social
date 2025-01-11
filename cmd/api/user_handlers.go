package api

import (
	"net/http"
	"strconv"

	"github.com/floroz/go-social/internal/domain"
	"github.com/go-chi/chi/v5"
)

func (app *Application) getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := getUserClaimFromContext(r.Context())
	if !ok {
		handleErrors(w, domain.NewUnauthorizedError("unauthorized"))
		return
	}

	user, err := app.UserService.GetByID(r.Context(), claims.ID)
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, user)
}

func (app *Application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "id"))
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

	user, err := app.UserService.Update(r.Context(), int64(userId), updateUserDto)
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, user)
}
