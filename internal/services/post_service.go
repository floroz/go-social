package services

import (
	"context"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
	"github.com/floroz/go-social/internal/validation"
)

type postService struct {
	postRepo interfaces.PostRepository
}

func NewPostService(postRepo interfaces.PostRepository) interfaces.PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) Create(ctx context.Context, createPost *domain.CreatePostDTO) (*domain.Post, error) {
	err := validation.ValidateCreatePostDTO(createPost)

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
