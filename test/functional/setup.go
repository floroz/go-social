package integration_tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/floroz/go-social/cmd/api"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/env"
	"github.com/floroz/go-social/internal/repositories"
	"github.com/floroz/go-social/internal/services"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const (
	signupEndpoint = "/api/v1/auth/signup"
	loginEndpoint  = "/api/v1/auth/login"
	logoutEndpoint = "/api/v1/auth/logout"
)

func runMigrations(db *sql.DB, migrationsDir string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	return nil
}

func startAPIServer(db *sql.DB) func() {
	env.MustLoadEnv("../../.env.local")

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	commentRepo := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo)

	postRepo := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepo, commentRepo)

	authService := services.NewAuthService(userRepo)

	config := &api.Config{
		Port: env.GetEnvValue("PORT"),
	}

	app := &api.Application{
		Config:         config,
		UserService:    userService,
		PostService:    postService,
		CommentService: commentService,
		AuthService:    authService,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.Config.Port),
		Handler:      app.Routes(),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Minute,
	}

	log.Info().Msgf("Started server on %s", app.Config.Port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("server error")
		}
	}()

	time.Sleep(2 * time.Second)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("server shutdown error")
		}
		log.Info().Msg("Server gracefully stopped")
	}
}

// signupAndGetCookies signs up a user and returns the cookies
func signupAndGetCookies(t *testing.T, client *http.Client, baseURL string, createUserDTO *domain.CreateUserDTO) (*domain.User, []*http.Cookie) {
	body, err := json.Marshal(createUserDTO)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, baseURL+signupEndpoint, bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var cookieNames []string
	for _, cookie := range resp.Cookies() {
		if cookie.Value != "" {
			cookieNames = append(cookieNames, cookie.Name)
		}
	}

	// status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	// cookies
	assert.NotEmpty(t, resp.Cookies())
	assert.Len(t, resp.Cookies(), 2)
	assert.Contains(t, cookieNames, "refresh_token")
	assert.Contains(t, cookieNames, "access_token")

	var signupResponse struct {
		Data domain.User `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&signupResponse)
	assert.NoError(t, err)

	return &signupResponse.Data, resp.Cookies()
}
