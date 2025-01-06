package services

import (
	"context"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
)

type postService struct {
	postRepo interfaces.PostRepository
}

func NewPostService(postRepo interfaces.PostRepository) interfaces.PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) Create(ctx context.Context, createPost *domain.CreatePostDTO) (*domain.Post, error) {
	return s.postRepo.Create(ctx, createPost)
}
