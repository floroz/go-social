package api

import (
	"net/http"
	"strconv"

	"github.com/floroz/go-social/internal/domain"
)

func (app *Application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := getUserClaimFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	createPostDTO := &domain.CreatePostDTO{}

	if err := readJSON(r.Body, createPostDTO); err != nil {
		handleErrors(w, domain.NewBadRequestError("failed to read request body"))
		return
	}

	post, err := app.PostService.Create(r.Context(), claims.ID, createPostDTO)
	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, post)
}

func (app *Application) listPostsHandler(w http.ResponseWriter, r *http.Request) {

	posts, err := app.PostService.List(r.Context(), 10, 0)

	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, posts)
}

func (app *Application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	claims, ok := getUserClaimFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleErrors(w, domain.NewBadRequestError("invalid id"))
		return
	}

	updatePostDTO := &domain.UpdatePostDTO{}

	if err := readJSON(r.Body, updatePostDTO); err != nil {
		handleErrors(w, domain.NewInternalServerError("failed to read request body"))
		return
	}

	post, err := app.PostService.Update(r.Context(), claims.ID, int64(postId), updatePostDTO)

	if err != nil {
		handleErrors(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, post)

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

	writeJSONResponse(w, http.StatusOK, post)
}
