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
	postRepo interfaces.PostRepository
}

func NewPostService(postRepo interfaces.PostRepository) interfaces.PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) Create(ctx context.Context, createPost *domain.CreatePostDTO) (*domain.Post, error) {
	err := validation.Validate.Struct(createPost)

	if err != nil {
		return nil, domain.NewBadRequestError(err.Error())
	}

	// validate that the user is trying to create a post for themselves
	// TODO: requires authentication

	post, err := s.postRepo.Create(ctx, createPost)

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

func (r *postService) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	// TODO: Who can request users posts? Are they all public, or users can decide whether to make them public or not?

	post, err := r.postRepo.GetByID(ctx, id)

	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return nil, domain.NewNotFoundError("post not found")
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to get post by id")
		return nil, domain.NewInternalServerError("failed to get post by id")
	}

	return post, nil
}
