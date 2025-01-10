package middlewares

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/env"
)

type contextKey string

const (
	ContextKeyUser contextKey = "user"
)

func AuthMiddleware(next http.Handler) http.Handler {
	jwtSecret := env.GetJWTSecret()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")

		if err != nil {
			http.Error(w, "missing access token", http.StatusUnauthorized)
			return
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &domain.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*domain.UserClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), ContextKeyUser, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
		}
	})
}
