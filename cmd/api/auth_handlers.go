package api

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/env"
)

const (
	accessTokenMaxDuration  = 15 * time.Minute
	refreshTokenMaxDuration = 24 * time.Hour
)

func (app *Application) signupHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: how do we protect this endpoint? This is a public endpoint, but we should consider rate limiting to protect against DDoS attacks

	createUserDto := &domain.CreateUserDTO{}

	if err := readJSON(r.Body, createUserDto); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := app.UserService.Create(r.Context(), createUserDto)
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
	loginUserDto := &domain.LoginUserDTO{}

	if err := readJSON(r.Body, loginUserDto); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := app.AuthService.Login(r.Context(), loginUserDto)
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
