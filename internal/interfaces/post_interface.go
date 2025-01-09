package interfaces

import (
	"context"

	"github.com/floroz/go-social/internal/domain"
)

type PostRepository interface {
	Create(ctx context.Context, post *domain.CreatePostDTO) (*domain.Post, error)
	List(ctx context.Context, limit int, offset int) ([]*domain.Post, error)
	GetByID(ctx context.Context, id int) (*domain.Post, error)
	Update(ctx context.Context, post *domain.UpdatePostDTO) (*domain.Post, error)
	Delete(ctx context.Context, id int) error
}

type PostService interface {
	Create(ctx context.Context, createPost *domain.CreatePostDTO) (*domain.Post, error)
	List(ctx context.Context, limit int, offset int) ([]*domain.Post, error)
	GetByID(ctx context.Context, id int) (*domain.Post, error)
	Update(ctx context.Context, post *domain.UpdatePostDTO) (*domain.Post, error)
	Delete(ctx context.Context, id int) error
}
