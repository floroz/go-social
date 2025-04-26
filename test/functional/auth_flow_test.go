package integration_tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/env"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.3",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Run migrations
	migrationsDir := "../../cmd/migrate/migrations"
	if err := runMigrations(db, migrationsDir); err != nil {
		log.Fatalf("Could not run migrations: %s", err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	shutdown := startAPIServer(db)
	defer shutdown()

	// run tests
	m.Run()
}

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, env.GetEnvValue("API_URL")+"/api/v1/healthz", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestUserSignup(t *testing.T) {
	// Generate unique suffix for email and username
	uniqueSuffix := fmt.Sprintf("ts%d", time.Now().UnixNano()) // Removed underscore
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "John",
			LastName:  "Doe",
			Email:     fmt.Sprintf("john.doe%s@example.com", uniqueSuffix), // Email can have special chars, keep suffix as is for uniqueness
			Username:  fmt.Sprintf("johndoets%d", time.Now().UnixNano()),   // Ensure username is purely alphanumeric
		},
		Password: "password123",
	}

	_, err := json.Marshal(createUserDTO)
	assert.NoError(t, err)

	user, _ := signupAndGetCookies(t, &http.Client{}, env.GetEnvValue("API_URL"), createUserDTO)

	// Assert
	assert.Equal(t, createUserDTO.Email, user.Email)
	assert.Equal(t, createUserDTO.Username, user.Username)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
	assert.Empty(t, user.Password)
}

func TestUserLogin(t *testing.T) {
	// Generate unique suffix for email and username
	uniqueSuffix := fmt.Sprintf("ts%d", time.Now().UnixNano()) // Removed underscore
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "John",
			LastName:  "Doe",
			Email:     fmt.Sprintf("john.doe%s@example.com", uniqueSuffix), // Email can have special chars, keep suffix as is for uniqueness
			Username:  fmt.Sprintf("johndoets%d", time.Now().UnixNano()),   // Ensure username is purely alphanumeric
		},
		Password: "password123",
	}

	_, cookies := signupAndGetCookies(t, &http.Client{}, env.GetEnvValue("API_URL"), createUserDTO)

	loginUserDTO := &domain.LoginUserDTO{
		Email:    createUserDTO.Email,
		Password: createUserDTO.Password,
	}

	body, err := json.Marshal(map[string]any{
		"data": loginUserDTO,
	})
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, env.GetEnvValue("API_URL")+loginEndpoint, bytes.NewBuffer(body))

	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookies[0])
	req.AddCookie(cookies[1])

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Define a struct matching the response structure {"data": User}
	var responseData struct {
		Data domain.User `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err)

	// Assert against the user data inside the "data" key
	userResponse := responseData.Data
	assert.Equal(t, createUserDTO.Email, userResponse.Email)
	assert.Equal(t, createUserDTO.Username, userResponse.Username)
	assert.NotZero(t, userResponse.ID)
	assert.NotZero(t, userResponse.CreatedAt)
	assert.NotZero(t, userResponse.UpdatedAt)
	assert.NotZero(t, userResponse.LastLogin)
	assert.Empty(t, userResponse.Password)

}
