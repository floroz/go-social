package api

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/floroz/go-social/internal/apitypes"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/env"
	"github.com/floroz/go-social/internal/errorcodes"
	"github.com/rs/zerolog/log"
)

const (
	accessTokenMaxDuration  = 15 * time.Minute
	refreshTokenMaxDuration = 24 * time.Hour
)

func (app *Application) signupHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: how do we protect this endpoint? This is a public endpoint, but we should consider rate limiting to protect against DDoS attacks

	var requestBody struct {
		Data *domain.CreateUserDTO `json:"data"`
	}

	if err := readJSON(r.Body, &requestBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error(), errorcodes.CodeBadRequest)
		return
	}

	user, err := app.UserService.Create(r.Context(), requestBody.Data)
	if err != nil {
		handleErrors(w, err)
		return
	}

	accessToken, err := app.AuthService.GenerateJWTToken(user, accessTokenMaxDuration)
	if err != nil {
		handleErrors(w, err)
		return
	}

	refreshToken, err := app.AuthService.GenerateJWTToken(user, refreshTokenMaxDuration)
	if err != nil {
		handleErrors(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Expires:  time.Now().Add(15 * time.Minute),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Wrap the user object in the standardized response structure
	// Map domain.User to apitypes.User
	apiUser := apitypes.User{
		Id:        &user.ID, // Use address-of for pointer type
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     apitypes.Email(user.Email), // Cast to apitypes.Email
		CreatedAt: &user.CreatedAt,            // Use address-of for pointer type
		UpdatedAt: &user.UpdatedAt,            // Use address-of for pointer type
		// Map LastLogin carefully, handling potential nil pointer if domain.User.LastLogin is *time.Time
		// Assuming apitypes.User.LastLogin is *time.Time based on common generation patterns
		// If domain.User.LastLogin is *time.Time and generated.User.LastLogin is *string:
		// LastLogin: func() *string {
		// 	if user.LastLogin != nil {
		// 		t := user.LastLogin.Format(time.RFC3339)
		// 		return &t
		// 	}
		// 	return nil
		// }(),
		// If both are time.Time or *time.Time, direct assignment might work,
		// but check generated type definition. Let's assume direct mapping for now if types match.
		// LastLogin: user.LastLogin, // Adjust based on actual generated type
	}
	// Handle potential nil LastLogin in domain.User if necessary
	if user.LastLogin != nil {
		// Assuming generated.User.LastLogin is *time.Time or time.Time
		// Adjust the assignment based on the exact type in generated.User
		// If domain.User.LastLogin is *time.Time and apitypes.User.LastLogin is *time.Time:
		apiUser.LastLogin = user.LastLogin // Correct assignment
	}
	// Note: The above assumes apitypes.User.LastLogin is *time.Time.
	// If it's different (e.g., *string), the mapping needs adjustment.

	response := apitypes.SignupSuccessResponse{
		Data: apiUser,
	}
	writeJSONResponse(w, http.StatusCreated, response)

}

func (app *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Data *domain.LoginUserDTO `json:"data"`
	}

	if err := readJSON(r.Body, &requestBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error(), errorcodes.CodeBadRequest)
		return
	}

	user, err := app.AuthService.Login(r.Context(), requestBody.Data)
	if err != nil {
		handleErrors(w, err)
		return
	}

	accessToken, err := app.AuthService.GenerateJWTToken(user, accessTokenMaxDuration)
	if err != nil {
		handleErrors(w, err)
		return
	}

	refreshToken, err := app.AuthService.GenerateJWTToken(user, refreshTokenMaxDuration)
	if err != nil {
		handleErrors(w, err)
		return
	}

	err = app.UserService.UpdateLastLogin(r.Context(), user.ID)

	if err != nil {
		// Log the error but continue, as failing to update last_login shouldn't block login
		log.Error().Err(err).Msg("failed to update last login")
	} else {
		// Re-fetch user data to get the updated LastLogin timestamp
		updatedUser, fetchErr := app.UserService.GetByID(r.Context(), user.ID)
		if fetchErr != nil {
			// Log error but proceed with the original user data if re-fetch fails
			log.Error().Err(fetchErr).Msg("failed to re-fetch user after updating last login")
		} else {
			user = updatedUser // Use the updated user data for the response
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Expires:  time.Now().Add(15 * time.Minute),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	})

	// Construct the API response containing the token
	apiLoginResponse := apitypes.LoginResponse{
		Token: accessToken, // Use the generated access token
	}

	// Wrap the login response in the standard success wrapper
	response := apitypes.LoginSuccessResponse{
		Data: apiLoginResponse,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

func (app *Application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
}

func (app *Application) refreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "missing refresh token", errorcodes.CodeUnauthorized)
		return
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &domain.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(env.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		writeJSONError(w, http.StatusUnauthorized, "invalid refresh token", errorcodes.CodeUnauthorized)
		return
	}

	claims, ok := token.Claims.(*domain.UserClaims)

	if !ok || !token.Valid {
		writeJSONError(w, http.StatusUnauthorized, "invalid refresh token claims", errorcodes.CodeUnauthorized)
		return
	}

	user := &domain.User{
		ID:        claims.ID,
		Username:  claims.Username,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
		Email:     claims.Email,
		CreatedAt: claims.CreatedAt,
		UpdatedAt: claims.UpdatedAt,
	}

	accessToken, err := app.AuthService.GenerateJWTToken(user, accessTokenMaxDuration)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to generate access token", errorcodes.CodeInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Expires:  time.Now().Add(15 * time.Minute),
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
}
