package api

import (
	"net/http"

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

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to write response")
	}
}

func (app *Application) listPostsHandler(w http.ResponseWriter, r *http.Request) {

	posts, err := app.PostService.List(r.Context(), 10, 0)

	if err != nil {
		handleErrors(w, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, posts); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to write response")
	}
}

func (app *Application) deletePostHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) updatePostHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {

}
