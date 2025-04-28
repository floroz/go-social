package api

import (
	"net/http"
	"strconv"

	"github.com/floroz/go-social/internal/apitypes"
	"github.com/floroz/go-social/internal/domain"
	// Error codes used via handleErrors in api.go
)

func (app *Application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("postId"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid post id"))
		return
	}

	userClaim, ok := getUserClaimFromContext(r.Context())
	if !ok {
		handleErrors(w, domain.NewUnauthorizedError("unauthorized"))
		return
	}
	userId := userClaim.ID // Get user ID after checking claim

	// Read request body into API type
	apiRequest := &apitypes.CreateCommentRequest{}
	if err := readJSON(r.Body, apiRequest); err != nil {
		handleErrors(w, domain.NewBadRequestError("failed to read request body: "+err.Error()))
		return
	}

	// Correctly map API request to domain DTO using embedded struct initialization
	domainDTO := &domain.CreateCommentDTO{
		EditableCommentFields: domain.EditableCommentFields{
			Content: apiRequest.Content,
		},
	}

	// Call service
	comment, err := app.CommentService.Create(r.Context(), userId, int64(postId), domainDTO)
	if err != nil {
		handleErrors(w, err)
		return
	}

	// Map domain.Comment to apitypes.Comment
	apiComment := mapDomainToApiComment(comment) // Use helper function

	// Wrap in success response
	response := apitypes.CreateCommentSuccessResponse{
		Data: apiComment,
	}

	writeJSONResponse(w, http.StatusCreated, response)
}

// Helper function to map domain.Comment to apitypes.Comment
func mapDomainToApiComment(comment *domain.Comment) apitypes.Comment {
	apiComment := apitypes.Comment{
		Id:        &comment.ID,     // Pointer
		PostId:    &comment.PostID, // Pointer, assuming generated type uses PostId
		UserId:    &comment.UserID, // Pointer, assuming generated type uses UserId
		Content:   comment.Content,
		CreatedAt: &comment.CreatedAt, // Pointer
		UpdatedAt: &comment.UpdatedAt, // Pointer
	}
	// Add mapping for other fields if they exist in apitypes.Comment
	return apiComment
}

// Helper function to map slice of domain.Comment to slice of apitypes.Comment
func mapDomainToApiComments(comments []domain.Comment) []apitypes.Comment {
	apiComments := make([]apitypes.Comment, len(comments))
	for i, c := range comments {
		apiComments[i] = mapDomainToApiComment(&c)
	}
	return apiComments
}

func (app *Application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid comment id"))
		return
	}

	postId, err := strconv.Atoi(r.PathValue("postId"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid post id"))
		return
	}

	userClaim, ok := getUserClaimFromContext(r.Context())
	if !ok {
		handleErrors(w, domain.NewUnauthorizedError("unauthorized"))
		return
	}
	userId := userClaim.ID // Get user ID after check

	// Read request body into API type
	apiRequest := &apitypes.UpdateCommentRequest{}
	if err := readJSON(r.Body, apiRequest); err != nil {
		handleErrors(w, domain.NewBadRequestError("failed to read request body: "+err.Error()))
		return
	}

	// Map API request to domain DTO
	domainDTO := &domain.UpdateCommentDTO{
		ID: int64(commentId), // ID comes from path param, not body
		EditableCommentFields: domain.EditableCommentFields{
			Content: apiRequest.Content,
		},
	}

	// Call service
	comment, err := app.CommentService.Update(r.Context(), userId, int64(postId), int64(commentId), domainDTO)
	if err != nil {
		handleErrors(w, err)
		return
	}

	// Map domain.Comment to apitypes.Comment
	apiComment := mapDomainToApiComment(comment)

	// Wrap in success response
	response := apitypes.UpdateCommentSuccessResponse{
		Data: apiComment,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

func (app *Application) getCommentByIdHandler(w http.ResponseWriter, r *http.Request) {
	commentId, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	comment, err := app.CommentService.GetByID(r.Context(), int64(commentId))

	if err != nil {
		handleErrors(w, err)
		return
	}

	// Map domain.Comment to apitypes.Comment
	apiComment := mapDomainToApiComment(comment)

	// Wrap in success response
	response := apitypes.GetCommentSuccessResponse{
		Data: apiComment,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

func (app *Application) listByPostIdHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("postId"))

	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid post id"))
		return
	}

	limitQuery := r.URL.Query().Get("limit")
	offsetQuery := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitQuery)

	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetQuery)

	if err != nil {
		offset = 0
	}

	if limitQuery == "" {
		limit = 10
	}

	comments, err := app.CommentService.ListByPostID(r.Context(), int64(postId), limit, offset)

	if err != nil {
		handleErrors(w, err)
		return
	}

	// Map domain comments to API comments
	apiComments := mapDomainToApiComments(comments)

	// Wrap in success response
	response := apitypes.ListCommentsSuccessResponse{
		Data: apiComments,
		// Add metadata here if implementing pagination
	}

	writeJSONResponse(w, http.StatusOK, response)
}

func (app *Application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid comment id"))
		return
	}

	userClaim, ok := getUserClaimFromContext(r.Context())
	if !ok {
		handleErrors(w, domain.NewUnauthorizedError("unauthorized"))
		return
	}

	err = app.CommentService.Delete(r.Context(), userClaim.ID, int64(commentId))
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSONResponse(w, http.StatusNoContent, nil)
}
