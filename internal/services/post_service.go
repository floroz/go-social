package services

import (
	"context"
	"errors"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
	"github.com/floroz/go-social/internal/validation"
	"github.com/rs/zerolog/log"
)

type postService struct {
	postRepo    interfaces.PostRepository
	commentRepo interfaces.CommentRepository
}

func NewPostService(postRepo interfaces.PostRepository, commentRepo interfaces.CommentRepository) interfaces.PostService {
	return &postService{postRepo: postRepo, commentRepo: commentRepo}
}

func (s *postService) Create(ctx context.Context, userId int64, createPost *domain.CreatePostDTO) (*domain.Post, error) {
	err := validation.Validate.Struct(createPost)

	if err != nil {
		return nil, domain.NewBadRequestError(err.Error())
	}

	post, err := s.postRepo.Create(ctx, userId, createPost)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *postService) List(ctx context.Context, limit int, offset int) ([]*domain.Post, error) {
	if limit > 100 {
		log.Warn().Msg("limit is too high, setting to 100")
		limit = 100
	} else if limit == 0 {
		limit = 10
	}

	posts, err := r.postRepo.List(ctx, limit, offset)

	if err != nil {
		log.Error().Err(err).Msg("failed to list posts")
		return nil, domain.NewInternalServerError("failed to list posts")
	}

	return posts, nil
}

func (r *postService) GetByID(ctx context.Context, postId int64) (*domain.Post, error) {
	// TODO: Who can request users posts? Are they all public, or users can decide whether to make them public or not?
	post, err := r.postRepo.GetByID(ctx, postId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.NewNotFoundError("post not found")
		}
		log.Error().Err(err).Msg("failed to get post by id")
		return nil, domain.NewInternalServerError("failed to get post by id")
	}

	// For now offset and limit are hard-coded
	comments, err := r.commentRepo.ListByPostID(ctx, postId, 100, 0)
	if err != nil {
		log.Error().Err(err).Msg("failed to get comments by post id")
		return nil, domain.NewInternalServerError("failed to get comments by post id")
	}

	post.Comments = comments

	return post, nil
}

func (r *postService) Update(ctx context.Context, userId, postId int64, updatedPost *domain.UpdatePostDTO) (*domain.Post, error) {
	if err := validation.Validate.Struct(updatedPost); err != nil {
		return nil, domain.NewBadRequestError(err.Error())
	}

	switch existingPost, err := r.postRepo.GetByID(ctx, postId); {
	case err != nil && err == domain.ErrNotFound:
		return nil, domain.NewNotFoundError("comment not found")
	case err != nil:
		return nil, domain.NewInternalServerError("failed to delete comment")
	case existingPost.UserID != userId:
		return nil, domain.NewForbiddenError("not allowed to delete comment")
	}

	post, err := r.postRepo.Update(ctx, userId, postId, updatedPost)

	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return nil, domain.NewNotFoundError("post not found")
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to update post")
		return nil, domain.NewInternalServerError("failed to update post")
	}

	return post, nil
}

func (r *postService) Delete(ctx context.Context, userId, postId int64) error {
	post, err := r.postRepo.GetByID(ctx, postId)

	switch {
	case err != nil && errors.Is(err, domain.ErrNotFound):
		return domain.NewNotFoundError("post not found")
	case err != nil:
		log.Error().Err(err).Msg("failed to get post by id")
		return domain.NewInternalServerError("failed to get post by id")
	case post.UserID != userId:
		return domain.NewForbiddenError("not allowed to delete post")
	}

	err = r.postRepo.Delete(ctx, userId, postId)

	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return domain.NewNotFoundError("post not found")
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to delete post")
		return domain.NewInternalServerError("failed to delete post")
	}

	return nil
}
