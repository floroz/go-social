package integration_tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/floroz/go-social/cmd/api"
	"github.com/floroz/go-social/internal/apitypes"
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
	signupEndpoint  = "/api/v1/auth/signup"
	loginEndpoint   = "/api/v1/auth/login"
	logoutEndpoint  = "/api/v1/auth/logout"
	healthzEndpoint = "/api/healthz"
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

func startTestAPIServer(db *sql.DB) *httptest.Server {
	env.MustLoadEnv("../../.env.local")

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	commentRepo := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo)

	postRepo := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepo, commentRepo)

	authService := services.NewAuthService(userRepo)

	// Create the application instance (Config is not strictly needed by httptest)
	// If your app initialization *requires* config, load it here.
	// For now, assuming it's not critical for route setup.
	app := &api.Application{
		// Config:         &api.Config{}, // Pass empty or load if needed
		UserService:    userService,
		PostService:    postService,
		CommentService: commentService,
		AuthService:    authService,
	}

	testServer := httptest.NewServer(app.Routes())

	log.Info().Msgf("Started test server on %s", testServer.URL)

	return testServer
}

// signupAndGetCookies signs up a user and returns the user data (as apitypes.User) and cookies
// It now uses the baseURL provided by the httptest server
func signupAndGetCookies(t *testing.T, client *http.Client, baseURL string, createUserDTO *domain.CreateUserDTO) (*apitypes.User, []*http.Cookie) {
	body, err := json.Marshal(map[string]any{
		"data": createUserDTO,
	})
	assert.NoError(t, err)

	// Construct full URL using baseURL and the specific endpoint constant
	req, err := http.NewRequest(http.MethodPost, baseURL+signupEndpoint, bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Use the provided client (which might be testServer.Client() or a custom one)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var cookieNames []string
	for _, cookie := range resp.Cookies() {
		if cookie.Value != "" {
			cookieNames = append(cookieNames, cookie.Name)
		}
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status 201 Created for signup")
	assert.NotEmpty(t, resp.Cookies(), "Expected cookies to be set on signup")
	if assert.Len(t, resp.Cookies(), 2, "Expected 2 cookies (access & refresh)") {
		assert.Contains(t, cookieNames, "refresh_token", "Expected refresh_token cookie")
		assert.Contains(t, cookieNames, "access_token", "Expected access_token cookie")
	}

	// Decode into the API response structure which contains apitypes.User
	var signupResponse struct {
		Data apitypes.User `json:"data"` // Use apitypes.User here
	}

	err = json.NewDecoder(resp.Body).Decode(&signupResponse)
	assert.NoError(t, err, "Failed to decode signup response body")

	return &signupResponse.Data, resp.Cookies()
}
