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

}

func (app *Application) updatePostHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {

}
