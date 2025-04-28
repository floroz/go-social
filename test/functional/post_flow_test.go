package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/floroz/go-social/internal/apitypes"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/errorcodes"
	"github.com/stretchr/testify/assert"
)

const (
	postsEndpoint = "/api/v1/posts"
)

func TestCreatePost(t *testing.T) {
	// Arrange: Sign up a user first
	// Use timestamp for unique suffix to ensure valid email/username
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Post",
			LastName:  "User",
			Email:     fmt.Sprintf("post.user%s@example.com", uniqueSuffix),
			Username:  fmt.Sprintf("postuser%s", uniqueSuffix),
		},
		Password: "password123",
	}
	signedUpUser, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotNil(t, signedUpUser)
	assert.NotEmpty(t, cookies)

	// Prepare create post request
	createPostReq := apitypes.CreatePostRequest{
		Content: "My first post content!",
	}
	// Marshal the request struct directly, as per the OpenAPI spec
	body, err := json.Marshal(createPostReq)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, testServerURL+postsEndpoint, bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	// Add auth cookies
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	client := testServer.Client()

	// Act: Perform create post request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert: Check status code and response body
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status 201 Created for post creation")

	var responseData apitypes.CreatePostSuccessResponse // Expect the wrapped response
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err, "Failed to decode create post response body")

	// Assert nested data
	createdPost := responseData.Data
	assert.NotNil(t, createdPost.Id, "Created Post ID should not be nil")
	if createdPost.Id != nil {
		assert.NotZero(t, *createdPost.Id, "Created Post ID should not be zero")
	}
	assert.NotNil(t, createdPost.UserId, "Created Post UserID should not be nil")
	if createdPost.UserId != nil {
		assert.Equal(t, *signedUpUser.Id, *createdPost.UserId, "Created Post UserID should match signed up user ID")
	}
	assert.Equal(t, createPostReq.Content, createdPost.Content)
	assert.NotNil(t, createdPost.CreatedAt, "CreatedAt should not be nil")
	if createdPost.CreatedAt != nil {
		assert.False(t, (*createdPost.CreatedAt).IsZero(), "CreatedAt should not be zero time")
	}
	assert.NotNil(t, createdPost.UpdatedAt, "UpdatedAt should not be nil")
	if createdPost.UpdatedAt != nil {
		assert.False(t, (*createdPost.UpdatedAt).IsZero(), "UpdatedAt should not be zero time")
	}
}

func TestGetPost(t *testing.T) {
	// Arrange: Sign up user and create a post
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "GetPost", LastName: "User",
			Email:    fmt.Sprintf("getpost.user%s@example.com", uniqueSuffix),
			Username: fmt.Sprintf("getpostuser%s", uniqueSuffix),
		},
		Password: "password123",
	}
	signedUpUser, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotNil(t, signedUpUser)
	assert.NotEmpty(t, cookies)

	createPostReq := apitypes.CreatePostRequest{Content: "Content for GetPost test"}
	createBody, _ := json.Marshal(createPostReq)
	createReq, _ := http.NewRequest(http.MethodPost, testServerURL+postsEndpoint, bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		createReq.AddCookie(cookie)
	}
	client := testServer.Client()
	createResp, _ := client.Do(createReq)
	assert.Equal(t, http.StatusCreated, createResp.StatusCode)
	var createRespData apitypes.CreatePostSuccessResponse
	_ = json.NewDecoder(createResp.Body).Decode(&createRespData)
	createResp.Body.Close()
	createdPostId := *createRespData.Data.Id // Get the ID of the created post

	// --- Test Case 1: Get existing post ---
	getReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s/%d", testServerURL, postsEndpoint, createdPostId), nil)
	assert.NoError(t, err)
	for _, cookie := range cookies {
		getReq.AddCookie(cookie)
	} // Add auth cookies

	// Act: Perform get post request
	getResp, err := client.Do(getReq)
	assert.NoError(t, err)
	defer getResp.Body.Close()

	// Assert: Check status code and response body
	assert.Equal(t, http.StatusOK, getResp.StatusCode, "Expected status 200 OK for get post")

	var getRespData apitypes.GetPostSuccessResponse // Expect the wrapped response
	err = json.NewDecoder(getResp.Body).Decode(&getRespData)
	assert.NoError(t, err, "Failed to decode get post response body")

	// Assert nested data
	retrievedPost := getRespData.Data
	assert.Equal(t, createdPostId, *retrievedPost.Id)
	assert.Equal(t, createPostReq.Content, retrievedPost.Content)
	assert.Equal(t, *signedUpUser.Id, *retrievedPost.UserId)

	// --- Test Case 2: Get non-existent post ---
	nonExistentId := int64(999999)
	getNonExistentReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s/%d", testServerURL, postsEndpoint, nonExistentId), nil)
	assert.NoError(t, err)
	for _, cookie := range cookies {
		getNonExistentReq.AddCookie(cookie)
	}

	// Act: Perform get non-existent post request
	getNonExistentResp, err := client.Do(getNonExistentReq)
	assert.NoError(t, err)
	defer getNonExistentResp.Body.Close()

	// Assert: Check status code (expect 404 Not Found)
	assert.Equal(t, http.StatusNotFound, getNonExistentResp.StatusCode, "Expected status 404 Not Found for non-existent post")

	// Assert error response structure
	var errorRespData apitypes.ApiErrorResponse
	err = json.NewDecoder(getNonExistentResp.Body).Decode(&errorRespData)
	assert.NoError(t, err, "Failed to decode error response body for non-existent post")
	assert.NotEmpty(t, errorRespData.Errors, "Expected errors array in response")
	if len(errorRespData.Errors) > 0 {
		// Compare the string value of the constant with the string field from the response
		assert.Equal(t, string(errorcodes.CodeNotFound), errorRespData.Errors[0].Code)
	}
}

func TestListPosts(t *testing.T) {
	// Arrange: Sign up user and create a couple of posts
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	// Ensure username is within DB limits (e.g., varchar(30))
	// Use a shorter prefix and potentially truncate timestamp if needed
	shortSuffix := uniqueSuffix[len(uniqueSuffix)-10:] // Example: take last 10 digits
	username := fmt.Sprintf("lstusr%s", shortSuffix)
	if len(username) > 30 { // Ensure it fits if suffix is long
		username = username[:30]
	}
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "ListPost", LastName: "User",
			Email:    fmt.Sprintf("listpost.user%s@example.com", uniqueSuffix), // Email length usually less restrictive
			Username: username,
		},
		Password: "password123",
	}
	_, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotEmpty(t, cookies)
	client := testServer.Client()

	// Create Post 1
	createPostReq1 := apitypes.CreatePostRequest{Content: "List Post 1"}
	body1, _ := json.Marshal(createPostReq1)
	req1, _ := http.NewRequest(http.MethodPost, testServerURL+postsEndpoint, bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req1.AddCookie(cookie)
	}
	resp1, _ := client.Do(req1)
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)
	resp1.Body.Close()

	// Create Post 2
	createPostReq2 := apitypes.CreatePostRequest{Content: "List Post 2"}
	body2, _ := json.Marshal(createPostReq2)
	req2, _ := http.NewRequest(http.MethodPost, testServerURL+postsEndpoint, bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req2.AddCookie(cookie)
	}
	resp2, _ := client.Do(req2)
	assert.Equal(t, http.StatusCreated, resp2.StatusCode)
	resp2.Body.Close()

	// --- Test Case: List posts ---
	listReq, err := http.NewRequest(http.MethodGet, testServerURL+postsEndpoint, nil)
	assert.NoError(t, err)
	for _, cookie := range cookies {
		listReq.AddCookie(cookie)
	} // Add auth cookies

	// Act: Perform list posts request
	listResp, err := client.Do(listReq)
	assert.NoError(t, err)
	defer listResp.Body.Close()

	// Assert: Check status code and response body
	assert.Equal(t, http.StatusOK, listResp.StatusCode, "Expected status 200 OK for list posts")

	var listRespData apitypes.ListPostsSuccessResponse // Expect the wrapped response
	err = json.NewDecoder(listResp.Body).Decode(&listRespData)
	assert.NoError(t, err, "Failed to decode list posts response body")

	// Assert nested data
	assert.NotEmpty(t, listRespData.Data, "Expected posts array in response data")
	// Check if at least the two created posts are present (order might vary)
	found1 := false
	found2 := false
	for _, post := range listRespData.Data {
		if post.Content == createPostReq1.Content {
			found1 = true
		}
		if post.Content == createPostReq2.Content {
			found2 = true
		}
	}
	assert.True(t, found1, "Expected to find 'List Post 1' in the list")
	assert.True(t, found2, "Expected to find 'List Post 2' in the list")
	// Note: More robust checks could involve checking IDs if predictable or fetching posts individually
}

func TestUpdatePost(t *testing.T) {
	// Arrange: Sign up user and create a post
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	// Ensure username is within DB limits
	shortSuffix := uniqueSuffix[len(uniqueSuffix)-10:]
	username := fmt.Sprintf("updusr%s", shortSuffix)
	if len(username) > 30 {
		username = username[:30]
	}
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "UpdatePost", LastName: "User",
			Email:    fmt.Sprintf("updatepost.user%s@example.com", uniqueSuffix),
			Username: username,
		},
		Password: "password123",
	}
	signedUpUser, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotNil(t, signedUpUser)
	assert.NotEmpty(t, cookies)
	client := testServer.Client()

	// Create the initial post
	createPostReq := apitypes.CreatePostRequest{Content: "Initial content for update"}
	createBody, _ := json.Marshal(createPostReq)
	createReq, _ := http.NewRequest(http.MethodPost, testServerURL+postsEndpoint, bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		createReq.AddCookie(cookie)
	}
	createResp, _ := client.Do(createReq)
	assert.Equal(t, http.StatusCreated, createResp.StatusCode)
	var createRespData apitypes.CreatePostSuccessResponse
	_ = json.NewDecoder(createResp.Body).Decode(&createRespData)
	createResp.Body.Close()
	createdPostId := *createRespData.Data.Id

	// --- Test Case 1: Successful update ---
	updatePostReq := apitypes.UpdatePostRequest{Content: "Updated post content!"}
	updateBody, err := json.Marshal(updatePostReq) // Update request body doesn't need 'data' wrapper
	assert.NoError(t, err)
	updateUrl := fmt.Sprintf("%s%s/%d", testServerURL, postsEndpoint, createdPostId)
	updateReq, err := http.NewRequest(http.MethodPut, updateUrl, bytes.NewBuffer(updateBody))
	assert.NoError(t, err)
	updateReq.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		updateReq.AddCookie(cookie)
	}

	// Act: Perform update post request
	updateResp, err := client.Do(updateReq)
	assert.NoError(t, err)
	defer updateResp.Body.Close()

	// Assert: Check status code and response body
	assert.Equal(t, http.StatusOK, updateResp.StatusCode, "Expected status 200 OK for update post")

	var updateRespData apitypes.UpdatePostSuccessResponse // Expect the wrapped response
	err = json.NewDecoder(updateResp.Body).Decode(&updateRespData)
	assert.NoError(t, err, "Failed to decode update post response body")

	// Assert nested data
	updatedPost := updateRespData.Data
	assert.Equal(t, createdPostId, *updatedPost.Id)
	assert.Equal(t, updatePostReq.Content, updatedPost.Content) // Check for updated content
	assert.Equal(t, *signedUpUser.Id, *updatedPost.UserId)
	// Optionally check if UpdatedAt changed
	assert.NotEqual(t, *createRespData.Data.UpdatedAt, *updatedPost.UpdatedAt)

	// --- Test Case 2: Update non-existent post ---
	nonExistentId := int64(999999)
	updateNonExistentReq := apitypes.UpdatePostRequest{Content: "Trying to update non-existent"}
	updateNonExistentBody, _ := json.Marshal(updateNonExistentReq)
	updateNonExistentUrl := fmt.Sprintf("%s%s/%d", testServerURL, postsEndpoint, nonExistentId)
	updateNonExistentHttpReq, _ := http.NewRequest(http.MethodPut, updateNonExistentUrl, bytes.NewBuffer(updateNonExistentBody))
	updateNonExistentHttpReq.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		updateNonExistentHttpReq.AddCookie(cookie)
	}

	// Act & Assert
	updateNonExistentResp, _ := client.Do(updateNonExistentHttpReq)
	defer updateNonExistentResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, updateNonExistentResp.StatusCode, "Expected 404 for updating non-existent post")

	// --- Test Case 3: Update another user's post (Forbidden) ---
	// Arrange: Sign up another user
	otherShortSuffix := uniqueSuffix[len(uniqueSuffix)-10:] // Use different suffix part if needed, or same is fine
	otherUsername := fmt.Sprintf("other%s", otherShortSuffix)
	if len(otherUsername) > 30 {
		otherUsername = otherUsername[:30]
	}
	otherUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Other", LastName: "User",
			Email:    fmt.Sprintf("other%s@example.com", uniqueSuffix),
			Username: otherUsername,
		},
		Password: "password123",
	}
	_, otherCookies := signupAndGetCookies(t, client, testServerURL, otherUserDTO)

	// Act: Try to update the first user's post with the second user's cookies
	updateOtherReq := apitypes.UpdatePostRequest{Content: "Trying to update other's post"}
	updateOtherBody, _ := json.Marshal(updateOtherReq)
	updateOtherHttpReq, _ := http.NewRequest(http.MethodPut, updateUrl, bytes.NewBuffer(updateOtherBody)) // Use original post URL
	updateOtherHttpReq.Header.Set("Content-Type", "application/json")
	for _, cookie := range otherCookies {
		updateOtherHttpReq.AddCookie(cookie)
	} // Use OTHER user's cookies

	// Assert
	updateOtherResp, _ := client.Do(updateOtherHttpReq)
	defer updateOtherResp.Body.Close()
	assert.Equal(t, http.StatusForbidden, updateOtherResp.StatusCode, "Expected 403 Forbidden for updating another user's post")

}

func TestDeletePost(t *testing.T) {
	// Arrange: Sign up user and create a post
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	// Ensure username is within DB limits
	shortSuffix := uniqueSuffix[len(uniqueSuffix)-10:]
	username := fmt.Sprintf("delusr%s", shortSuffix)
	if len(username) > 30 {
		username = username[:30]
	}
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "DeletePost", LastName: "User",
			Email:    fmt.Sprintf("deletepost.user%s@example.com", uniqueSuffix),
			Username: username,
		},
		Password: "password123",
	}
	signedUpUser, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotNil(t, signedUpUser)
	assert.NotEmpty(t, cookies)
	client := testServer.Client()

	// Create the post to be deleted
	createPostReq := apitypes.CreatePostRequest{Content: "Post to be deleted"}
	createBody, _ := json.Marshal(createPostReq)
	createReq, _ := http.NewRequest(http.MethodPost, testServerURL+postsEndpoint, bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		createReq.AddCookie(cookie)
	}
	createResp, _ := client.Do(createReq)
	assert.Equal(t, http.StatusCreated, createResp.StatusCode)
	var createRespData apitypes.CreatePostSuccessResponse
	_ = json.NewDecoder(createResp.Body).Decode(&createRespData)
	createResp.Body.Close()
	createdPostId := *createRespData.Data.Id

	// --- Test Case 1: Successful delete ---
	deleteUrl := fmt.Sprintf("%s%s/%d", testServerURL, postsEndpoint, createdPostId)
	deleteReq, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	assert.NoError(t, err)
	for _, cookie := range cookies {
		deleteReq.AddCookie(cookie)
	}

	// Act: Perform delete post request
	deleteResp, err := client.Do(deleteReq)
	assert.NoError(t, err)
	defer deleteResp.Body.Close()

	// Assert: Check status code (204 No Content)
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode, "Expected status 204 No Content for delete post")

	// Assert: Verify post is actually deleted (GET should return 404)
	getReq, _ := http.NewRequest(http.MethodGet, deleteUrl, nil)
	for _, cookie := range cookies {
		getReq.AddCookie(cookie)
	}
	getResp, _ := client.Do(getReq)
	defer getResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode, "Expected status 404 Not Found after deleting post")

	// --- Test Case 2: Delete non-existent post ---
	nonExistentId := int64(999999)
	deleteNonExistentUrl := fmt.Sprintf("%s%s/%d", testServerURL, postsEndpoint, nonExistentId)
	deleteNonExistentReq, _ := http.NewRequest(http.MethodDelete, deleteNonExistentUrl, nil)
	for _, cookie := range cookies {
		deleteNonExistentReq.AddCookie(cookie)
	}

	// Act & Assert
	deleteNonExistentResp, _ := client.Do(deleteNonExistentReq)
	defer deleteNonExistentResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, deleteNonExistentResp.StatusCode, "Expected 404 for deleting non-existent post")

	// --- Test Case 3: Delete another user's post (Forbidden) ---
	// Arrange: Create another post by the first user
	createPostReq2 := apitypes.CreatePostRequest{Content: "Another post"}
	createBody2, _ := json.Marshal(createPostReq2)
	createReq2, _ := http.NewRequest(http.MethodPost, testServerURL+postsEndpoint, bytes.NewBuffer(createBody2))
	createReq2.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		createReq2.AddCookie(cookie)
	}
	createResp2, _ := client.Do(createReq2)
	assert.Equal(t, http.StatusCreated, createResp2.StatusCode)
	var createRespData2 apitypes.CreatePostSuccessResponse
	_ = json.NewDecoder(createResp2.Body).Decode(&createRespData2)
	createResp2.Body.Close()
	postToKeepId := *createRespData2.Data.Id

	// Arrange: Sign up another user
	otherShortSuffix := uniqueSuffix[len(uniqueSuffix)-10:]
	otherUsername := fmt.Sprintf("otherdel%s", otherShortSuffix)
	if len(otherUsername) > 30 {
		otherUsername = otherUsername[:30]
	}
	otherUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{FirstName: "OtherDel", LastName: "UserDel", Email: fmt.Sprintf("otherdel%s@example.com", uniqueSuffix), Username: otherUsername},
		Password:          "password123",
	}
	_, otherCookies := signupAndGetCookies(t, client, testServerURL, otherUserDTO)

	// Act: Try to delete the first user's post with the second user's cookies
	deleteOtherUrl := fmt.Sprintf("%s%s/%d", testServerURL, postsEndpoint, postToKeepId)
	deleteOtherHttpReq, _ := http.NewRequest(http.MethodDelete, deleteOtherUrl, nil)
	for _, cookie := range otherCookies {
		deleteOtherHttpReq.AddCookie(cookie)
	} // Use OTHER user's cookies

	// Assert
	deleteOtherResp, _ := client.Do(deleteOtherHttpReq)
	defer deleteOtherResp.Body.Close()
	assert.Equal(t, http.StatusForbidden, deleteOtherResp.StatusCode, "Expected 403 Forbidden for deleting another user's post")

}

// TODO: Add tests for validation errors, authorization errors (e.g., updating/deleting others' posts)
