package api

import (
	"net/http"

	"github.com/floroz/go-social/internal/apitypes"
	"github.com/floroz/go-social/internal/domain"
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

	// Map domain.User to apitypes.User
	apiUser := apitypes.User{
		Id:        &user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     apitypes.Email(user.Email),
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
		LastLogin: user.LastLogin, // Assuming types match (*time.Time)
	}

	// Wrap in the success response structure
	response := apitypes.GetUserProfileSuccessResponse{
		Data: apiUser,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

func (app *Application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context claims, not URL param
	claims, ok := getUserClaimFromContext(r.Context())
	if !ok {
		handleErrors(w, domain.NewUnauthorizedError("unauthorized"))
		return
	}
	userId := claims.ID // Use ID from claims

	updateUserDto := &domain.UpdateUserDTO{}

	var err error
	err = readJSON(r.Body, updateUserDto)
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("failed to read request body: "+err.Error()))
		return
	}

	// Call service using the correct userId type (int64)
	user, err := app.UserService.Update(r.Context(), userId, updateUserDto)
	if err != nil {
		handleErrors(w, err)
		return
	}

	// Map domain.User to apitypes.User
	apiUser := apitypes.User{
		Id:        &user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     apitypes.Email(user.Email),
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
		LastLogin: user.LastLogin, // Assuming types match (*time.Time)
	}

	// Wrap in the success response structure
	response := apitypes.UpdateUserProfileSuccessResponse{
		Data: apiUser,
	}

	writeJSONResponse(w, http.StatusOK, response)
}
