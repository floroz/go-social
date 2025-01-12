package api

import (
	"net/http"
	"strconv"

	"github.com/floroz/go-social/internal/domain"
)

func (app *Application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("postId"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	userClaim, ok := getUserClaimFromContext(r.Context())
	userId := userClaim.ID
	if !ok {
		handleErrors(w, domain.NewBadRequestError("invalid user claim"))
		return
	}

	createCommentDTO := &domain.CreateCommentDTO{}
	if err := readJSON(r.Body, createCommentDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := app.CommentService.Create(r.Context(), userId, int64(postId), createCommentDTO)
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, comment)
}

func (app *Application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	postId, err := strconv.Atoi(r.PathValue("postId"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid post id"))
		return
	}

	userClaim, ok := getUserClaimFromContext(r.Context())
	userId := userClaim.ID
	if !ok {
		handleErrors(w, domain.NewBadRequestError("invalid user claim"))
		return
	}

	updateCommentDTO := &domain.UpdateCommentDTO{
		ID: int64(commentId),
	}
	if err := readJSON(r.Body, updateCommentDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := app.CommentService.Update(r.Context(), userId, int64(postId), int64(commentId), updateCommentDTO)
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, comment)
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

	writeJSONResponse(w, http.StatusOK, comment)
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

	writeJSONResponse(w, http.StatusOK, comments)
}

func (app *Application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	userClaim, ok := getUserClaimFromContext(r.Context())
	if !ok {
		handleErrors(w, domain.NewBadRequestError("invalid user claim"))
		return
	}

	err = app.CommentService.Delete(r.Context(), userClaim.ID, int64(commentId))
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSONResponse(w, http.StatusNoContent, nil)
}
