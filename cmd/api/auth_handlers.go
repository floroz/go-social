package api

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/env"
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
		writeJSONError(w, http.StatusBadRequest, err.Error())
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

	writeJSONResponse(w, http.StatusCreated, user)

}

func (app *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Data *domain.LoginUserDTO `json:"data"`
	}

	if err := readJSON(r.Body, &requestBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
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

	writeJSONResponse(w, http.StatusOK, user)
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
		writeJSONError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &domain.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(env.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		writeJSONError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	claims, ok := token.Claims.(*domain.UserClaims)

	if !ok || !token.Valid {
		writeJSONError(w, http.StatusUnauthorized, "invalid refresh token claims")
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
		writeJSONError(w, http.StatusInternalServerError, "failed to generate access token")
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
