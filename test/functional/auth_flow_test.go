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

	"github.com/floroz/go-social/internal/apitypes"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/errorcodes"
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
	uniqueSuffix := fmt.Sprintf("ts%d", time.Now().UnixMilli()) // Use Milli for shorter suffix
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "John",
			LastName:  "Doe",
			Email:     fmt.Sprintf("john.doe.%s@example.com", uniqueSuffix), // Adjusted format
			Username:  fmt.Sprintf("johndoe%s", uniqueSuffix),               // Adjusted format
		},
		Password: "password123",
	}
	_, err := json.Marshal(createUserDTO)
	assert.NoError(t, err)

	// Act: Perform signup request directly
	requestPayload := map[string]any{"data": createUserDTO}
	bodyBytes, err := json.Marshal(requestPayload)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, testServerURL+signupEndpoint, bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := testServer.Client()
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Assert response body structure {"data": User}
	var signupSuccessResponse struct {
		Data apitypes.User `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&signupSuccessResponse)
	assert.NoError(t, err, "Failed to decode signup success response")
	assert.NotNil(t, signupSuccessResponse.Data, "Response 'data' field should not be nil")

	// Assert fields of the user within the data wrapper
	userFromResponse := signupSuccessResponse.Data
	assert.Equal(t, createUserDTO.Email, string(userFromResponse.Email))
	assert.Equal(t, createUserDTO.Username, userFromResponse.Username)
	assert.NotNil(t, userFromResponse.Id, "User ID should not be nil")
	if userFromResponse.Id != nil {
		assert.NotZero(t, *userFromResponse.Id, "User ID should not be zero")
	}
	assert.NotNil(t, userFromResponse.CreatedAt, "CreatedAt should not be nil")
	if userFromResponse.CreatedAt != nil {
		assert.False(t, (*userFromResponse.CreatedAt).IsZero(), "CreatedAt should not be zero time")
	}
	assert.NotNil(t, userFromResponse.UpdatedAt, "UpdatedAt should not be nil")
	if userFromResponse.UpdatedAt != nil {
		assert.False(t, (*userFromResponse.UpdatedAt).IsZero(), "UpdatedAt should not be zero time")
	}
	// Password field is not expected in the response.

	// We can still use signupAndGetCookies for setting up other tests if needed,
	// but TestUserSignup itself now does direct validation.
	// For example, to get cookies for TestUserLogin:
	// _, cookiesForLogin := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO) // This would be a separate user for isolation
	// Note: The actual signupAndGetCookies is assumed to be in setup.go or similar.
}

func TestUserSignup_DuplicateEmail(t *testing.T) {
	// Arrange: Sign up an initial user
	// Use the existing signupAndGetCookies from setup (assuming it's available and works with wrapped requests)
	uniqueSuffix1 := fmt.Sprintf("ts%d", time.Now().UnixMilli()) // Use Milli
	initialUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Initial", LastName: "User",
			Email:    fmt.Sprintf("dup.email.%s@example.com", uniqueSuffix1), // Adjusted
			Username: fmt.Sprintf("inituser%s", uniqueSuffix1),               // Adjusted
		},
		Password: "password123",
	}
	_, _ = signupAndGetCookies(t, testServer.Client(), testServerURL, initialUserDTO) // Sign up the first user

	// Attempt to sign up another user with the same email
	uniqueSuffix2 := fmt.Sprintf("ts%d", time.Now().UnixMilli()) // Use Milli for potentially different parts
	duplicateEmailUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Duplicate", LastName: "Email",
			Email:    initialUserDTO.Email,                     // Same email
			Username: fmt.Sprintf("dupemail%s", uniqueSuffix2), // Different username
		},
		Password: "password456",
	}
	requestPayload := map[string]any{"data": duplicateEmailUserDTO}
	bodyBytes, err := json.Marshal(requestPayload)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, testServerURL+signupEndpoint, bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := testServer.Client()
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert: Check status code and error response structure
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var errorResponse apitypes.ApiErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.NoError(t, err, "Failed to decode error response for duplicate email")
	assert.NotEmpty(t, errorResponse.Errors, "Errors array should not be empty")
	if len(errorResponse.Errors) > 0 {
		assert.Equal(t, string(errorcodes.CodeConflict), errorResponse.Errors[0].Code) // Cast to string for comparison
	}
}

func TestUserSignup_DuplicateUsername(t *testing.T) {
	// Arrange: Sign up an initial user
	uniqueSuffix1 := fmt.Sprintf("ts%d", time.Now().UnixMilli()) // Use Milli
	initialUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Initial", LastName: "User",
			Email:    fmt.Sprintf("init.user.%s@example.com", uniqueSuffix1), // Adjusted
			Username: fmt.Sprintf("dupuser%s", uniqueSuffix1),                // Adjusted
		},
		Password: "password123",
	}
	_, _ = signupAndGetCookies(t, testServer.Client(), testServerURL, initialUserDTO)

	// Attempt to sign up another user with the same username
	uniqueSuffix2 := fmt.Sprintf("ts%d", time.Now().UnixMilli()) // Use Milli
	duplicateUsernameUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Duplicate", LastName: "Username",
			Email:    fmt.Sprintf("another%s@example.com", uniqueSuffix2), // Different email
			Username: initialUserDTO.Username,                             // Same username
		},
		Password: "password456",
	}
	requestPayload := map[string]any{"data": duplicateUsernameUserDTO}
	bodyBytes, err := json.Marshal(requestPayload)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, testServerURL+signupEndpoint, bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := testServer.Client()
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert: Check status code and error response structure
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var errorResponse apitypes.ApiErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.NoError(t, err, "Failed to decode error response for duplicate username")
	assert.NotEmpty(t, errorResponse.Errors, "Errors array should not be empty")
	if len(errorResponse.Errors) > 0 {
		assert.Equal(t, string(errorcodes.CodeConflict), errorResponse.Errors[0].Code) // Cast to string for comparison
	}
}

func TestUserSignup_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name          string
		dto           *domain.CreateUserDTO
		expectedField string // Optional: for field-specific validation errors
		expectedCode  string // Ensure this is string type
	}{
		{
			name: "Password too short",
			dto: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "Valid", LastName: "User", Email: "valid@example.com", Username: "validuser1",
				},
				Password: "short",
			},
			expectedField: "password",
			expectedCode:  string(errorcodes.CodeValidationError), // Cast to string
		},
		{
			name: "Invalid email format",
			dto: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "Valid", LastName: "User", Email: "invalid-email", Username: "validuser2",
				},
				Password: "password123",
			},
			expectedField: "email",
			expectedCode:  string(errorcodes.CodeValidationError), // Cast to string
		},
		{
			name: "Missing first name",
			dto: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					LastName: "User", Email: "valid3@example.com", Username: "validuser3",
				},
				Password: "password123",
			},
			expectedField: "first_name",
			expectedCode:  string(errorcodes.CodeValidationError), // Cast to string
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Ensure unique email/username for fields not being tested for validation errors, to avoid unintended 409s
			// Use a simpler, shorter unique part for these test cases to avoid DB length issues for valid fields.
			// Also ensure username is alphanumeric.
			ts := fmt.Sprintf("%d", time.Now().UnixMilli()%10000) // Shorter unique part
			if tc.dto.Email == "valid@example.com" {
				tc.dto.Email = fmt.Sprintf("val%s@example.com", ts)
			}
			if tc.dto.Username == "validuser1" {
				tc.dto.Username = fmt.Sprintf("valusr%s", ts)
			}
			if tc.dto.Email == "invalid-email" && (tc.name == "Invalid email format") {
				// Keep it invalid for this specific test, but ensure other fields are unique if needed
				tc.dto.Username = fmt.Sprintf("valusr%s", ts)
			} else if tc.dto.Email == "valid3@example.com" { // For missing first name test
				tc.dto.Email = fmt.Sprintf("valfn%s@example.com", ts)
			}

			if tc.dto.Username == "validuser2" && (tc.name == "Invalid email format") {
				// Username is fine here, email is being tested
			} else if tc.dto.Username == "validuser3" && (tc.name == "Missing first name") {
				tc.dto.Username = fmt.Sprintf("valfnusr%s", ts)
			}

			requestPayload := map[string]any{"data": tc.dto}
			bodyBytes, err := json.Marshal(requestPayload)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, testServerURL+signupEndpoint, bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			client := testServer.Client()
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var errorResponse apitypes.ApiErrorResponse
			err = json.NewDecoder(resp.Body).Decode(&errorResponse)
			assert.NoError(t, err, "Failed to decode validation error response")
			assert.NotEmpty(t, errorResponse.Errors, "Errors array should not be empty for validation error")

			if len(errorResponse.Errors) > 0 {
				foundExpectedError := false
				for _, apiErr := range errorResponse.Errors {
					if tc.expectedCode != "" {
						assert.Equal(t, tc.expectedCode, apiErr.Code)
					}
					if tc.expectedField != "" && apiErr.Field != nil && *apiErr.Field == tc.expectedField {
						foundExpectedError = true
						// Optionally check message if it's stable
					} else if tc.expectedField == "" { // If no specific field, any validation error is fine
						foundExpectedError = true
					}
				}
				if tc.expectedField != "" { // Only assert found if a specific field was expected
					assert.True(t, foundExpectedError, "Expected error for field '%s' not found in errors: %+v", tc.expectedField, errorResponse.Errors)
				}
			}
		})
	}
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

	// Assert: Check status code and response body for the token
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var responseData struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err, "Failed to decode login response body")
	assert.NotEmpty(t, responseData.Data.Token, "Expected token in login response")
	// Optionally, add more checks for the token format if needed
}

func TestUserLogout(t *testing.T) {
	// Arrange: Sign up a user first to get credentials and cookies
	uniqueSuffix := fmt.Sprintf("ts%d", time.Now().UnixNano())
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Logout", LastName: "User",
			Email:    fmt.Sprintf("logout.user%s@example.com", uniqueSuffix),
			Username: fmt.Sprintf("logoutuser%d", time.Now().UnixNano()),
		},
		Password: "password123",
	}
	_, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotEmpty(t, cookies, "Signup should provide cookies")

	// Prepare logout request
	req, err := http.NewRequest(http.MethodPost, testServerURL+logoutEndpoint, nil) // No body needed
	assert.NoError(t, err)
	// Add cookies obtained from signup
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	client := testServer.Client()

	// Act: Perform logout request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert: Check status code and cookies
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200 OK for logout")

	// Check that response cookies are expired/cleared
	foundExpiredAccess := false
	foundExpiredRefresh := false
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "access_token" {
			assert.True(t, cookie.Expires.Before(time.Now()), "Access token cookie should be expired")
			foundExpiredAccess = true
		}
		if cookie.Name == "refresh_token" {
			assert.True(t, cookie.Expires.Before(time.Now()), "Refresh token cookie should be expired")
			foundExpiredRefresh = true
		}
	}
	assert.True(t, foundExpiredAccess, "Expected expired access_token cookie in response")
	assert.True(t, foundExpiredRefresh, "Expected expired refresh_token cookie in response")

	// Assert: Attempt an authenticated request (e.g., get user profile) - should fail
	// Note: Need to define userProfileEndpoint or use an existing authenticated endpoint
	// For now, let's assume /api/v1/users/ is the profile endpoint
	profileReq, err := http.NewRequest(http.MethodGet, testServerURL+"/api/v1/users", nil)
	assert.NoError(t, err)
	// Intentionally DO NOT add cookies
	profileResp, err := client.Do(profileReq)
	assert.NoError(t, err)
	defer profileResp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, profileResp.StatusCode, "Expected 401 Unauthorized after logout")
}

func TestRefreshToken(t *testing.T) {
	// Arrange: Sign up a user first to get credentials and cookies
	uniqueSuffix := fmt.Sprintf("ts%d", time.Now().UnixNano())
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Refresh", LastName: "User",
			Email:    fmt.Sprintf("refresh.user%s@example.com", uniqueSuffix),
			Username: fmt.Sprintf("refreshuser%d", time.Now().UnixNano()),
		},
		Password: "password123",
	}
	_, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotEmpty(t, cookies, "Signup should provide cookies")

	// Find the refresh token cookie
	var refreshTokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "refresh_token" {
			refreshTokenCookie = c
			break
		}
	}
	assert.NotNil(t, refreshTokenCookie, "Refresh token cookie not found after signup")

	// Prepare refresh request
	refreshEndpoint := "/api/v1/auth/refresh"                                        // Define endpoint if not already done
	req, err := http.NewRequest(http.MethodPost, testServerURL+refreshEndpoint, nil) // No body needed
	assert.NoError(t, err)
	req.AddCookie(refreshTokenCookie) // Only send the refresh token cookie
	client := testServer.Client()

	// Act: Perform refresh request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert: Check status code and cookies
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200 OK for refresh")

	// Check that a new access_token cookie is set and valid (e.g., not expired immediately)
	foundNewAccessToken := false
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "access_token" {
			assert.NotEmpty(t, cookie.Value, "New access token cookie should not be empty")
			assert.True(t, cookie.Expires.After(time.Now()), "New access token cookie should not be immediately expired")
			// Optionally: compare value to old access token to ensure it's different
			foundNewAccessToken = true
			break // Found it
		}
	}
	assert.True(t, foundNewAccessToken, "Expected new access_token cookie in refresh response")

	// --- Test Case 2: Refresh with invalid/missing token ---
	reqInvalid, err := http.NewRequest(http.MethodPost, testServerURL+refreshEndpoint, nil)
	assert.NoError(t, err)
	// Intentionally DO NOT add refresh token cookie

	// Act & Assert
	respInvalid, err := client.Do(reqInvalid)
	assert.NoError(t, err)
	defer respInvalid.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, respInvalid.StatusCode, "Expected 401 Unauthorized for refresh without token")
}

func TestGetUserProfile(t *testing.T) {
	// Arrange: Sign up a user first
	uniqueSuffix := fmt.Sprintf("ts%d", time.Now().UnixNano())
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Profile", LastName: "User",
			Email:    fmt.Sprintf("profile.user%s@example.com", uniqueSuffix),
			Username: fmt.Sprintf("profileuser%d", time.Now().UnixNano()),
		},
		Password: "password123",
	}
	signedUpUser, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotNil(t, signedUpUser)
	assert.NotEmpty(t, cookies)

	// Prepare get profile request
	profileEndpoint := "/api/v1/users" // Define endpoint
	req, err := http.NewRequest(http.MethodGet, testServerURL+profileEndpoint, nil)
	assert.NoError(t, err)
	// Add cookies obtained from signup
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	client := testServer.Client()

	// Act: Perform get profile request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert: Check status code and response body
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200 OK for get profile")

	var responseData apitypes.GetUserProfileSuccessResponse // Expect the wrapped response
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err, "Failed to decode get profile response body")

	// Assert nested data matches the signed-up user
	profileUser := responseData.Data
	assert.Equal(t, *signedUpUser.Id, *profileUser.Id)
	assert.Equal(t, createUserDTO.FirstName, profileUser.FirstName)
	assert.Equal(t, createUserDTO.LastName, profileUser.LastName)
	assert.Equal(t, createUserDTO.Username, profileUser.Username)
	assert.Equal(t, createUserDTO.Email, string(profileUser.Email))
	assert.NotNil(t, profileUser.CreatedAt)
	assert.NotNil(t, profileUser.UpdatedAt)
	// LastLogin might be nil initially depending on flow, check if needed

	// --- Test Case 2: Get profile without authentication ---
	reqUnauth, err := http.NewRequest(http.MethodGet, testServerURL+profileEndpoint, nil)
	assert.NoError(t, err)
	// Intentionally DO NOT add cookies

	// Act & Assert
	respUnauth, err := client.Do(reqUnauth)
	assert.NoError(t, err)
	defer respUnauth.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, respUnauth.StatusCode, "Expected 401 Unauthorized for get profile without auth")
}

func TestUpdateUserProfile(t *testing.T) {
	// Arrange: Sign up a user first
	uniqueSuffix := fmt.Sprintf("ts%d", time.Now().UnixNano())
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Update", LastName: "User", // Use valid LastName
			Email:    fmt.Sprintf("updateme%s@example.com", uniqueSuffix),
			Username: fmt.Sprintf("updateme%d", time.Now().UnixNano()),
		},
		Password: "password123",
	}
	signedUpUser, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotNil(t, signedUpUser)
	assert.NotEmpty(t, cookies)

	// Prepare update profile request (PUT requires all fields)
	updateEndpoint := "/api/v1/users" // Define endpoint
	// Helper function to get pointer to string
	stringPtr := func(s string) *string { return &s }
	// Create apitypes.Email value and take its address
	emailVal := apitypes.Email(createUserDTO.Email)
	updateReqDTO := apitypes.UpdateUserProfileRequest{
		FirstName: stringPtr("UpdatedFirstName"),     // Update First Name
		LastName:  stringPtr("UpdatedLastName"),      // Update Last Name
		Email:     &emailVal,                         // Correct type: *apitypes.Email
		Username:  stringPtr(createUserDTO.Username), // Keep original Username
	}
	// Wrap the DTO in a "data" field for the request payload
	requestPayload := map[string]any{"data": updateReqDTO}
	body, err := json.Marshal(requestPayload)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, testServerURL+updateEndpoint, bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	// Add cookies obtained from signup
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	client := testServer.Client()

	// Act: Perform update profile request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert: Check status code and response body
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200 OK for update profile")

	var responseData apitypes.UpdateUserProfileSuccessResponse // Expect the wrapped response
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err, "Failed to decode update profile response body")

	// Assert nested data matches the updated user
	updatedUser := responseData.Data
	assert.Equal(t, *signedUpUser.Id, *updatedUser.Id)
	assert.Equal(t, *updateReqDTO.FirstName, updatedUser.FirstName)         // Check updated field (dereference)
	assert.Equal(t, *updateReqDTO.LastName, updatedUser.LastName)           // Check updated field (dereference)
	assert.Equal(t, *updateReqDTO.Username, updatedUser.Username)           // Check field sent (dereference)
	assert.Equal(t, string(*updateReqDTO.Email), string(updatedUser.Email)) // Corrected: Compare string values
	assert.NotEqual(t, *signedUpUser.UpdatedAt, *updatedUser.UpdatedAt)     // UpdatedAt should change

	// --- Test Case 2: Update profile without authentication ---
	reqUnauth, err := http.NewRequest(http.MethodPut, testServerURL+updateEndpoint, bytes.NewBuffer(body)) // Reuse body
	assert.NoError(t, err)
	reqUnauth.Header.Set("Content-Type", "application/json")
	// Intentionally DO NOT add cookies

	// Act & Assert
	respUnauth, err := client.Do(reqUnauth)
	assert.NoError(t, err)
	defer respUnauth.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, respUnauth.StatusCode, "Expected 401 Unauthorized for update profile without auth")

	// TODO: Add test case for validation errors (e.g., invalid email format if updating email)
}
