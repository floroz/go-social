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

	createCommentDTO := &domain.CreateCommentDTO{
		PostID: postId,
	}

	if err := readJSON(r.Body, createCommentDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := app.CommentService.Create(r.Context(), createCommentDTO)
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, comment)
}

func (app *Application) getCommentByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	comment, err := app.CommentService.GetByID(r.Context(), id)

	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusOK, comment)
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

	comments, err := app.CommentService.ListByPostID(r.Context(), postId, limit, offset)

	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusOK, comments)
}

func (app *Application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	err = app.CommentService.Delete(r.Context(), id)

	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func (app *Application) deleteByPostIdCommentHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("postId"))

	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid post id"))
		return
	}

	err = app.CommentService.DeleteByPostID(r.Context(), postId)

	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
