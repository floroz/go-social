package integration_tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"net/http/httptest"

	"github.com/floroz/go-social/internal/domain"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB
var testServer *httptest.Server // Store the test server instance
var testServerURL string        // Store the test server URL

func TestMain(m *testing.M) {
	log.Println("Setting up test database via dockertest...")
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

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

	if err := resource.Expire(120); err != nil {
		log.Fatalf("Could not set resource expiration: %s", err)
	}
	log.Println("Docker container expiration set to 120 seconds.")

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker database: %s", err)
	}

	// Schedule cleanup for the container
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
		log.Println("Docker container purged.")
	}()

	// Run migrations
	migrationsDir := "../../cmd/migrate/migrations"
	if err := runMigrations(db, migrationsDir); err != nil {
		log.Fatalf("Could not run migrations: %s", err)
	}
	log.Println("Migrations applied successfully.")

	// Start the test API server (using the established db connection)
	testServer = startTestAPIServer(db)
	testServerURL = testServer.URL
	defer testServer.Close()

	exitCode := m.Run()

	// --- Teardown (handled by defers) ---

	// Use os.Exit to exit with the test result code
	os.Exit(exitCode)
}

func TestHealthCheck(t *testing.T) {
	// Use the dynamic URL from the test server started in TestMain
	req, err := http.NewRequest(http.MethodGet, testServerURL+healthzEndpoint, nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	// Use the test server's client for direct requests
	client := testServer.Client()
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserSignup(t *testing.T) {
	// Arrange: Prepare the signup request
	// Generate unique suffix for email and username
	uniqueSuffix := fmt.Sprintf("ts%d", time.Now().UnixNano())
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "John",
			LastName:  "Doe",
			Email:     fmt.Sprintf("john.doe%s@example.com", uniqueSuffix),
			Username:  fmt.Sprintf("johndoets%d", time.Now().UnixNano()),
		},
		Password: "password123",
	}
	_, err := json.Marshal(createUserDTO)
	assert.NoError(t, err)

	// Act: Perform signup request
	user, _ := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)

	// Assert
	assert.Equal(t, createUserDTO.Email, user.Email)
	assert.Equal(t, createUserDTO.Username, user.Username)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
	assert.Empty(t, user.Password) // Ensure password hash is not returned
}

func TestUserLogin(t *testing.T) {
	// Arrange: Sign up a user first to get credentials and cookies
	uniqueSuffix := fmt.Sprintf("ts%d", time.Now().UnixNano())
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "John",
			LastName:  "Doe",
			Email:     fmt.Sprintf("john.doe%s@example.com", uniqueSuffix),
			Username:  fmt.Sprintf("johndoets%d", time.Now().UnixNano()),
		},
		Password: "password123",
	}
	_, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)

	// Prepare login request
	loginUserDTO := &domain.LoginUserDTO{
		Email:    createUserDTO.Email,
		Password: createUserDTO.Password,
	}
	body, err := json.Marshal(map[string]any{"data": loginUserDTO})
	assert.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, testServerURL+loginEndpoint, bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	// Add cookies obtained from signup
	if len(cookies) >= 2 {
		req.AddCookie(cookies[0])
		req.AddCookie(cookies[1])
	}
	client := testServer.Client()

	// Act: Perform login request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert: Check status code and response body
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var responseData struct {
		Data domain.User `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err)
	userResponse := responseData.Data
	assert.Equal(t, createUserDTO.Email, userResponse.Email)
	assert.Equal(t, createUserDTO.Username, userResponse.Username)
	assert.NotZero(t, userResponse.ID)
	assert.NotZero(t, userResponse.CreatedAt)
	assert.NotZero(t, userResponse.UpdatedAt)
	assert.NotZero(t, userResponse.LastLogin)
	assert.Empty(t, userResponse.Password)
}
