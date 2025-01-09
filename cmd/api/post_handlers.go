package api

import (
	"net/http"
	"strconv"

	"github.com/floroz/go-social/internal/domain"
)

func (app *Application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	createPostDTO := &domain.CreatePostDTO{}

	if err := readJSON(r.Body, createPostDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := app.PostService.Create(r.Context(), createPostDTO)
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, post)
}

func (app *Application) listPostsHandler(w http.ResponseWriter, r *http.Request) {

	posts, err := app.PostService.List(r.Context(), 10, 0)

	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusOK, posts)
}

func (app *Application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	err = app.PostService.Delete(r.Context(), id)

	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)

}

func (app *Application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	updatePostDTO := &domain.UpdatePostDTO{}

	if err := readJSON(r.Body, updatePostDTO); err != nil {
		handleErrors(w, domain.NewInternalServerError("failed to read request body"))
		return
	}

	updatePostDTO.ID = id

	post, err := app.PostService.Update(r.Context(), updatePostDTO)

	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusOK, post)

}

func (app *Application) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	post, err := app.PostService.GetByID(r.Context(), id)

	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSON(w, http.StatusOK, post)
}
