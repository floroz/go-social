package api

import (
	"net/http"

	"github.com/floroz/go-social/internal/apitypes"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/errorcodes"
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

	var requestBody struct {
		Data *domain.UpdateUserDTO `json:"data"`
	}

	err := readJSON(r.Body, &requestBody)
	if err != nil {
		// Use a more specific error code if readJSON fails due to malformed JSON or validation issues
		// For now, assuming CodeBadRequest is acceptable for general parsing/unmarshaling issues.
		// If readJSON itself could return validator.ValidationErrors, that would be handled by handleErrors.
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error(), errorcodes.CodeBadRequest, "")
		return
	}

	// Call service using the correct userId type (int64) and the unwrapped DTO
	user, err := app.UserService.Update(r.Context(), userId, requestBody.Data)
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
