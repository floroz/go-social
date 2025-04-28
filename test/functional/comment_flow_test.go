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
	"github.com/floroz/go-social/internal/errorcodes" // Import errorcodes
	"github.com/stretchr/testify/assert"
)

// Helper function to create a post and return its ID (requires signed-up user cookies)
func createPostForTest(t *testing.T, client *http.Client, cookies []*http.Cookie, content string) int64 {
	createPostReq := apitypes.CreatePostRequest{Content: content}
	body, err := json.Marshal(createPostReq)
	assert.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, testServerURL+postsEndpoint, bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var respData apitypes.CreatePostSuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	assert.NoError(t, err)
	resp.Body.Close()
	assert.NotNil(t, respData.Data.Id)
	return *respData.Data.Id
}

// Helper function to create a comment and return its ID
func createCommentForTest(t *testing.T, client *http.Client, cookies []*http.Cookie, postId int64, content string) int64 {
	createCommentReq := apitypes.CreateCommentRequest{Content: content}
	body, err := json.Marshal(createCommentReq)
	assert.NoError(t, err)
	commentUrl := fmt.Sprintf("%s%s/%d/comments", testServerURL, postsEndpoint, postId)
	req, err := http.NewRequest(http.MethodPost, commentUrl, bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var respData apitypes.CreateCommentSuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	assert.NoError(t, err)
	resp.Body.Close()
	assert.NotNil(t, respData.Data.Id)
	return *respData.Data.Id
}

func TestCreateComment(t *testing.T) {
	// Arrange: Sign up user and create a post
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Comment", LastName: "User",
			Email:    fmt.Sprintf("comment.user%s@example.com", uniqueSuffix),
			Username: fmt.Sprintf("commentuser%s", uniqueSuffix),
		},
		Password: "password123",
	}
	signedUpUser, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotNil(t, signedUpUser)
	assert.NotEmpty(t, cookies)
	client := testServer.Client()

	postId := createPostForTest(t, client, cookies, "Post for commenting")

	// Prepare create comment request
	createCommentReq := apitypes.CreateCommentRequest{
		Content: "My first comment!",
	}
	body, err := json.Marshal(createCommentReq) // Request body is just the content
	assert.NoError(t, err)

	commentUrl := fmt.Sprintf("%s%s/%d/comments", testServerURL, postsEndpoint, postId)
	req, err := http.NewRequest(http.MethodPost, commentUrl, bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Act: Perform create comment request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert: Check status code and response body
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status 201 Created for comment creation")

	var responseData apitypes.CreateCommentSuccessResponse // Expect the wrapped response
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err, "Failed to decode create comment response body")

	// Assert nested data
	createdComment := responseData.Data
	assert.NotNil(t, createdComment.Id, "Created Comment ID should not be nil")
	if createdComment.Id != nil {
		assert.NotZero(t, *createdComment.Id, "Created Comment ID should not be zero")
	}
	assert.NotNil(t, createdComment.UserId, "Created Comment UserID should not be nil")
	if createdComment.UserId != nil {
		assert.Equal(t, *signedUpUser.Id, *createdComment.UserId, "Created Comment UserID should match signed up user ID")
	}
	assert.NotNil(t, createdComment.PostId, "Created Comment PostID should not be nil")
	if createdComment.PostId != nil {
		assert.Equal(t, postId, *createdComment.PostId, "Created Comment PostID should match the post ID")
	}
	assert.Equal(t, createCommentReq.Content, createdComment.Content)
	assert.NotNil(t, createdComment.CreatedAt, "CreatedAt should not be nil")
	assert.NotNil(t, createdComment.UpdatedAt, "UpdatedAt should not be nil")
}

func TestGetComment(t *testing.T) {
	// Arrange: Sign up user, create post, create comment
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	// Ensure username is within DB limits
	shortSuffix := uniqueSuffix[len(uniqueSuffix)-10:]
	username := fmt.Sprintf("getcmtusr%s", shortSuffix)
	if len(username) > 30 {
		username = username[:30]
	}
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "GetComment", LastName: "User",
			Email:    fmt.Sprintf("getcomment.user%s@example.com", uniqueSuffix),
			Username: username,
		},
		Password: "password123",
	}
	signedUpUser, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotNil(t, signedUpUser)
	assert.NotEmpty(t, cookies)
	client := testServer.Client()

	postId := createPostForTest(t, client, cookies, "Post for GetComment test")
	createdCommentId := createCommentForTest(t, client, cookies, postId, "Comment to get")

	// --- Test Case 1: Get existing comment ---
	getCommentUrl := fmt.Sprintf("%s%s/%d/comments/%d", testServerURL, postsEndpoint, postId, createdCommentId)
	getReq, err := http.NewRequest(http.MethodGet, getCommentUrl, nil)
	assert.NoError(t, err)
	for _, cookie := range cookies {
		getReq.AddCookie(cookie)
	} // Add auth cookies

	// Act: Perform get comment request
	getResp, err := client.Do(getReq)
	assert.NoError(t, err)
	defer getResp.Body.Close()

	// Assert: Check status code and response body
	assert.Equal(t, http.StatusOK, getResp.StatusCode, "Expected status 200 OK for get comment")

	var getRespData apitypes.GetCommentSuccessResponse // Expect the wrapped response
	err = json.NewDecoder(getResp.Body).Decode(&getRespData)
	assert.NoError(t, err, "Failed to decode get comment response body")

	// Assert nested data
	retrievedComment := getRespData.Data
	assert.Equal(t, createdCommentId, *retrievedComment.Id)
	// assert.Equal(t, createCommentReq.Content, retrievedComment.Content) // Content was defined in helper
	assert.Equal(t, *signedUpUser.Id, *retrievedComment.UserId)
	assert.Equal(t, postId, *retrievedComment.PostId)

	// --- Test Case 2: Get non-existent comment ---
	nonExistentId := int64(999999)
	getNonExistentUrl := fmt.Sprintf("%s%s/%d/comments/%d", testServerURL, postsEndpoint, postId, nonExistentId)
	getNonExistentReq, err := http.NewRequest(http.MethodGet, getNonExistentUrl, nil)
	assert.NoError(t, err)
	for _, cookie := range cookies {
		getNonExistentReq.AddCookie(cookie)
	}

	// Act & Assert
	getNonExistentResp, err := client.Do(getNonExistentReq)
	assert.NoError(t, err)
	defer getNonExistentResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, getNonExistentResp.StatusCode, "Expected status 404 Not Found for non-existent comment")
}

func TestListComments(t *testing.T) {
	// Arrange: Sign up user, create post, create comments
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	// Ensure username is within DB limits
	shortSuffix := uniqueSuffix[len(uniqueSuffix)-10:]
	username := fmt.Sprintf("lstcmtusr%s", shortSuffix)
	if len(username) > 30 {
		username = username[:30]
	}
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "ListComment", LastName: "User",
			Email:    fmt.Sprintf("listcomment.user%s@example.com", uniqueSuffix),
			Username: username,
		},
		Password: "password123",
	}
	_, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotEmpty(t, cookies)
	client := testServer.Client()

	postId := createPostForTest(t, client, cookies, "Post for ListComments test")

	// Create Comment 1 & 2 using helper
	comment1Content := "List Comment 1"
	comment2Content := "List Comment 2"
	_ = createCommentForTest(t, client, cookies, postId, comment1Content)
	_ = createCommentForTest(t, client, cookies, postId, comment2Content)

	// --- Test Case: List comments ---
	commentUrl := fmt.Sprintf("%s%s/%d/comments", testServerURL, postsEndpoint, postId)
	listReq, err := http.NewRequest(http.MethodGet, commentUrl, nil)
	assert.NoError(t, err)
	for _, cookie := range cookies {
		listReq.AddCookie(cookie)
	} // Add auth cookies

	// Act: Perform list comments request
	listResp, err := client.Do(listReq)
	assert.NoError(t, err)
	defer listResp.Body.Close()

	// Assert: Check status code and response body
	assert.Equal(t, http.StatusOK, listResp.StatusCode, "Expected status 200 OK for list comments")

	var listRespData apitypes.ListCommentsSuccessResponse // Expect the wrapped response
	err = json.NewDecoder(listResp.Body).Decode(&listRespData)
	assert.NoError(t, err, "Failed to decode list comments response body")

	// Assert nested data
	assert.NotEmpty(t, listRespData.Data, "Expected comments array in response data")
	// Check if at least the two created comments are present
	found1 := false
	found2 := false
	for _, comment := range listRespData.Data {
		if comment.Content == comment1Content {
			found1 = true
		}
		if comment.Content == comment2Content {
			found2 = true
		}
	}
	assert.True(t, found1, "Expected to find 'List Comment 1' in the list")
	assert.True(t, found2, "Expected to find 'List Comment 2' in the list")
}

func TestUpdateComment(t *testing.T) {
	// Arrange: Sign up user, create post, create comment
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	// Ensure username is within DB limits
	shortSuffix := uniqueSuffix[len(uniqueSuffix)-10:]
	username := fmt.Sprintf("updcmtusr%s", shortSuffix)
	if len(username) > 30 {
		username = username[:30]
	}
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "UpdateComment", LastName: "User",
			Email:    fmt.Sprintf("updatecomment.user%s@example.com", uniqueSuffix),
			Username: username,
		},
		Password: "password123",
	}
	signedUpUser, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotNil(t, signedUpUser)
	assert.NotEmpty(t, cookies)
	client := testServer.Client()

	postId := createPostForTest(t, client, cookies, "Post for UpdateComment test")
	createdCommentId := createCommentForTest(t, client, cookies, postId, "Initial comment content")

	// --- Test Case 1: Successful update ---
	updateCommentReq := apitypes.UpdateCommentRequest{Content: "Updated comment content!"}
	updateBody, err := json.Marshal(updateCommentReq)
	assert.NoError(t, err)
	updateUrl := fmt.Sprintf("%s%s/%d/comments/%d", testServerURL, postsEndpoint, postId, createdCommentId)
	updateReq, err := http.NewRequest(http.MethodPut, updateUrl, bytes.NewBuffer(updateBody))
	assert.NoError(t, err)
	updateReq.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		updateReq.AddCookie(cookie)
	}

	// Act: Perform update comment request
	updateResp, err := client.Do(updateReq)
	assert.NoError(t, err)
	defer updateResp.Body.Close()

	// Assert: Check status code and response body
	assert.Equal(t, http.StatusOK, updateResp.StatusCode, "Expected status 200 OK for update comment")

	var updateRespData apitypes.UpdateCommentSuccessResponse // Expect the wrapped response
	err = json.NewDecoder(updateResp.Body).Decode(&updateRespData)
	assert.NoError(t, err, "Failed to decode update comment response body")

	// Assert nested data
	updatedComment := updateRespData.Data
	assert.Equal(t, createdCommentId, *updatedComment.Id)
	assert.Equal(t, updateCommentReq.Content, updatedComment.Content) // Check for updated content
	assert.Equal(t, *signedUpUser.Id, *updatedComment.UserId)
	assert.Equal(t, postId, *updatedComment.PostId)

	// --- Test Case 2: Update non-existent comment ---
	nonExistentId := int64(999999)
	updateNonExistentReq := apitypes.UpdateCommentRequest{Content: "Trying to update non-existent"}
	updateNonExistentBody, _ := json.Marshal(updateNonExistentReq)
	updateNonExistentUrl := fmt.Sprintf("%s%s/%d/comments/%d", testServerURL, postsEndpoint, postId, nonExistentId)
	updateNonExistentHttpReq, _ := http.NewRequest(http.MethodPut, updateNonExistentUrl, bytes.NewBuffer(updateNonExistentBody))
	updateNonExistentHttpReq.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		updateNonExistentHttpReq.AddCookie(cookie)
	}

	// Act & Assert
	updateNonExistentResp, _ := client.Do(updateNonExistentHttpReq)
	defer updateNonExistentResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, updateNonExistentResp.StatusCode, "Expected 404 for updating non-existent comment")
	// Assert error response structure
	var errorRespDataNotFound apitypes.ApiErrorResponse
	err = json.NewDecoder(updateNonExistentResp.Body).Decode(&errorRespDataNotFound)
	assert.NoError(t, err, "Failed to decode error response body for non-existent comment update")
	assert.NotEmpty(t, errorRespDataNotFound.Errors, "Expected errors array in response")
	if len(errorRespDataNotFound.Errors) > 0 {
		assert.Equal(t, string(errorcodes.CodeNotFound), errorRespDataNotFound.Errors[0].Code) // Use errorcodes
	}

	// --- Test Case 3: Update another user's comment (Forbidden) ---
	// Arrange: Sign up another user
	otherShortSuffix := uniqueSuffix[len(uniqueSuffix)-10:] // Use different suffix part if needed
	otherUsername := fmt.Sprintf("otherupdcmt%s", otherShortSuffix)
	if len(otherUsername) > 30 {
		otherUsername = otherUsername[:30]
	}
	otherUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{FirstName: "OtherUpd", LastName: "UserUpd", Email: fmt.Sprintf("otherupd%s@example.com", uniqueSuffix), Username: otherUsername},
		Password:          "password123",
	}
	_, otherCookies := signupAndGetCookies(t, client, testServerURL, otherUserDTO)

	// Act: Try to update the first user's comment with the second user's cookies
	updateOtherReq := apitypes.UpdateCommentRequest{Content: "Trying to update other's comment"}
	updateOtherBody, _ := json.Marshal(updateOtherReq)
	updateOtherHttpReq, _ := http.NewRequest(http.MethodPut, updateUrl, bytes.NewBuffer(updateOtherBody)) // Use original comment URL
	updateOtherHttpReq.Header.Set("Content-Type", "application/json")
	for _, cookie := range otherCookies {
		updateOtherHttpReq.AddCookie(cookie)
	} // Use OTHER user's cookies

	// Assert
	updateOtherResp, _ := client.Do(updateOtherHttpReq)
	defer updateOtherResp.Body.Close()
	assert.Equal(t, http.StatusForbidden, updateOtherResp.StatusCode, "Expected 403 Forbidden for updating another user's comment")
	// Assert error response structure
	var errorRespDataForbidden apitypes.ApiErrorResponse
	err = json.NewDecoder(updateOtherResp.Body).Decode(&errorRespDataForbidden)
	assert.NoError(t, err, "Failed to decode error response body for forbidden comment update")
	assert.NotEmpty(t, errorRespDataForbidden.Errors, "Expected errors array in response")
	if len(errorRespDataForbidden.Errors) > 0 {
		assert.Equal(t, string(errorcodes.CodeForbidden), errorRespDataForbidden.Errors[0].Code) // Use errorcodes
	}
}

func TestDeleteComment(t *testing.T) {
	// Arrange: Sign up user, create post, create comment
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	shortSuffix := uniqueSuffix[len(uniqueSuffix)-10:]
	username := fmt.Sprintf("delcmtusr%s", shortSuffix)
	if len(username) > 30 {
		username = username[:30]
	}
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "DeleteComment", LastName: "User",
			Email:    fmt.Sprintf("deletecomment.user%s@example.com", uniqueSuffix),
			Username: username,
		},
		Password: "password123",
	}
	_, cookies := signupAndGetCookies(t, testServer.Client(), testServerURL, createUserDTO)
	assert.NotEmpty(t, cookies)
	client := testServer.Client()

	postId := createPostForTest(t, client, cookies, "Post for DeleteComment test")
	commentToDeleteId := createCommentForTest(t, client, cookies, postId, "Comment to delete")

	// --- Test Case 1: Successful delete ---
	deleteUrl := fmt.Sprintf("%s%s/%d/comments/%d", testServerURL, postsEndpoint, postId, commentToDeleteId)
	deleteReq, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	assert.NoError(t, err)
	for _, cookie := range cookies {
		deleteReq.AddCookie(cookie)
	}

	// Act: Perform delete comment request
	deleteResp, err := client.Do(deleteReq)
	assert.NoError(t, err)
	defer deleteResp.Body.Close()

	// Assert: Check status code (204 No Content)
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode, "Expected status 204 No Content for delete comment")

	// Assert: Try to GET the deleted comment, should be 404
	getDeletedReq, err := http.NewRequest(http.MethodGet, deleteUrl, nil)
	assert.NoError(t, err)
	for _, cookie := range cookies {
		getDeletedReq.AddCookie(cookie)
	}
	getDeletedResp, err := client.Do(getDeletedReq)
	assert.NoError(t, err)
	defer getDeletedResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, getDeletedResp.StatusCode, "Expected 404 Not Found when getting deleted comment")

	// --- Test Case 2: Delete non-existent comment ---
	nonExistentId := int64(999999)
	deleteNonExistentUrl := fmt.Sprintf("%s%s/%d/comments/%d", testServerURL, postsEndpoint, postId, nonExistentId)
	deleteNonExistentReq, _ := http.NewRequest(http.MethodDelete, deleteNonExistentUrl, nil)
	for _, cookie := range cookies {
		deleteNonExistentReq.AddCookie(cookie)
	}

	// Act & Assert
	deleteNonExistentResp, _ := client.Do(deleteNonExistentReq)
	defer deleteNonExistentResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, deleteNonExistentResp.StatusCode, "Expected 404 for deleting non-existent comment")
	// Assert error response structure
	var errorRespDataNotFound apitypes.ApiErrorResponse
	err = json.NewDecoder(deleteNonExistentResp.Body).Decode(&errorRespDataNotFound)
	assert.NoError(t, err, "Failed to decode error response body for non-existent comment delete")
	assert.NotEmpty(t, errorRespDataNotFound.Errors, "Expected errors array in response")
	if len(errorRespDataNotFound.Errors) > 0 {
		assert.Equal(t, string(errorcodes.CodeNotFound), errorRespDataNotFound.Errors[0].Code) // Use errorcodes
	}

	// --- Test Case 3: Delete another user's comment (Forbidden) ---
	// Arrange: Create another comment by the first user
	commentToDeleteByOtherId := createCommentForTest(t, client, cookies, postId, "Another comment to delete")
	// Sign up another user
	otherShortSuffix := uniqueSuffix[len(uniqueSuffix)-10:] + "b" // Ensure different suffix
	otherUsername := fmt.Sprintf("otherdelcmt%s", otherShortSuffix)
	if len(otherUsername) > 30 {
		otherUsername = otherUsername[:30]
	}
	otherUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{FirstName: "OtherDel", LastName: "UserDel", Email: fmt.Sprintf("otherdel%s@example.com", uniqueSuffix), Username: otherUsername},
		Password:          "password123",
	}
	_, otherCookies := signupAndGetCookies(t, client, testServerURL, otherUserDTO)

	// Act: Try to delete the first user's comment with the second user's cookies
	deleteOtherUrl := fmt.Sprintf("%s%s/%d/comments/%d", testServerURL, postsEndpoint, postId, commentToDeleteByOtherId)
	deleteOtherHttpReq, _ := http.NewRequest(http.MethodDelete, deleteOtherUrl, nil)
	for _, cookie := range otherCookies {
		deleteOtherHttpReq.AddCookie(cookie)
	} // Use OTHER user's cookies

	// Assert
	deleteOtherResp, _ := client.Do(deleteOtherHttpReq)
	defer deleteOtherResp.Body.Close()
	assert.Equal(t, http.StatusForbidden, deleteOtherResp.StatusCode, "Expected 403 Forbidden for deleting another user's comment")
	// Assert error response structure
	var errorRespDataForbidden apitypes.ApiErrorResponse
	err = json.NewDecoder(deleteOtherResp.Body).Decode(&errorRespDataForbidden)
	assert.NoError(t, err, "Failed to decode error response body for forbidden comment delete")
	assert.NotEmpty(t, errorRespDataForbidden.Errors, "Expected errors array in response")
	if len(errorRespDataForbidden.Errors) > 0 {
		assert.Equal(t, string(errorcodes.CodeForbidden), errorRespDataForbidden.Errors[0].Code) // Use errorcodes
	}
}

// TODO: Add tests for validation errors, authorization errors
