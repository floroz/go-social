package api

import (
	"net/http"
	"strconv"

	"github.com/floroz/go-social/internal/apitypes"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/errorcodes" // Added for CodeBadRequest
	"github.com/rs/zerolog/log"
)

func (app *Application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := getUserClaimFromContext(r.Context())
	if !ok {
		handleErrors(w, domain.NewUnauthorizedError("unauthorized"))
		return
	}

	var requestBody struct {
		Data *apitypes.CreatePostRequest `json:"data"`
	}
	if err := readJSON(r.Body, &requestBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error(), errorcodes.CodeBadRequest, "")
		return
	}

	// Correctly map API request to domain DTO using embedded struct initialization
	domainDTO := &domain.CreatePostDTO{
		EditablePostFields: domain.EditablePostFields{
			Content: requestBody.Data.Content,
		},
	}

	post, err := app.PostService.Create(r.Context(), claims.ID, domainDTO)
	if err != nil {
		handleErrors(w, err)
		return
	}

	// Map domain.Post to apitypes.Post
	apiPost := mapDomainToApiPost(post) // Use helper function (to be defined)

	// Wrap in success response
	response := apitypes.CreatePostSuccessResponse{
		Data: apiPost,
	}

	writeJSONResponse(w, http.StatusCreated, response)
}

// Helper function to map domain.Post to apitypes.Post
func mapDomainToApiPost(post *domain.Post) apitypes.Post {
	apiPost := apitypes.Post{
		Id:        &post.ID,     // Pointer
		UserId:    &post.UserID, // Corrected field name based on domain.Post
		Content:   post.Content,
		CreatedAt: &post.CreatedAt, // Pointer
		UpdatedAt: &post.UpdatedAt, // Pointer
	}
	// Add mapping for other fields if they exist in apitypes.Post (e.g., author username)
	return apiPost
}

// Helper function to map slice of domain.Post to slice of apitypes.Post
func mapDomainToApiPosts(posts []domain.Post) []apitypes.Post { // Accept []domain.Post
	apiPosts := make([]apitypes.Post, len(posts))
	for i, p := range posts {
		// Pass the address of the element p to mapDomainToApiPost which expects a pointer
		apiPosts[i] = mapDomainToApiPost(&p)
	}
	return apiPosts
}

func (app *Application) listPostsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Add pagination query parameter handling (page, limit)

	posts, err := app.PostService.List(r.Context(), 10, 0) // Using default limit for now

	if err != nil {
		handleErrors(w, err)
		return
	}

	// Map domain posts to API posts
	apiPosts := mapDomainToApiPosts(posts)

	// Wrap in success response
	response := apitypes.ListPostsSuccessResponse{
		Data: apiPosts,
		// Add metadata here if implementing pagination
	}

	writeJSONResponse(w, http.StatusOK, response)
}

func (app *Application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	claims, ok := getUserClaimFromContext(r.Context())
	if !ok {
		// Use standard error handling
		handleErrors(w, domain.NewUnauthorizedError("unauthorized"))
		return
	}

	err = app.PostService.Delete(r.Context(), claims.ID, int64(postId))
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSONResponse(w, http.StatusNoContent, nil)

}

func (app *Application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := getUserClaimFromContext(r.Context())
	if !ok {
		// Use standard error handling
		handleErrors(w, domain.NewUnauthorizedError("unauthorized"))
		return
	}

	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	var requestBody struct {
		Data *apitypes.UpdatePostRequest `json:"data"`
	}
	if err := readJSON(r.Body, &requestBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error(), errorcodes.CodeBadRequest, "")
		return
	}

	// Map API request to domain DTO
	domainDTO := &domain.UpdatePostDTO{
		EditablePostFields: domain.EditablePostFields{
			Content: requestBody.Data.Content,
		},
	}

	// Add logging before service call
	log.Debug().Int64("authUserID", claims.ID).Int("pathPostID", postId).Msg("Attempting to update post")

	// Call service
	post, err := app.PostService.Update(r.Context(), claims.ID, int64(postId), domainDTO)
	if err != nil {
		log.Error().Err(err).Int64("authUserID", claims.ID).Int("pathPostID", postId).Msg("Error calling PostService.Update")
		handleErrors(w, err)
		return
	}

	// Map domain.Post to apitypes.Post
	apiPost := mapDomainToApiPost(post)

	// Wrap in success response
	response := apitypes.UpdatePostSuccessResponse{
		Data: apiPost,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

func (app *Application) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	post, err := app.PostService.GetByID(r.Context(), int64(postId))

	if err != nil {
		handleErrors(w, err)
		return
	}

	// Map domain.Post to apitypes.Post
	apiPost := mapDomainToApiPost(post)

	// Wrap in success response
	response := apitypes.GetPostSuccessResponse{
		Data: apiPost,
	}

	writeJSONResponse(w, http.StatusOK, response)
}
