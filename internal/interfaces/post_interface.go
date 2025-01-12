package interfaces

import (
	"context"

	"github.com/floroz/go-social/internal/domain"
)

type PostRepository interface {
	Create(ctx context.Context, userId int64, post *domain.CreatePostDTO) (*domain.Post, error)
	List(ctx context.Context, limit int, offset int) ([]domain.Post, error)
	GetByID(ctx context.Context, postId int64) (*domain.Post, error)
	Update(ctx context.Context, userId, postId int64, post *domain.UpdatePostDTO) (*domain.Post, error)
	Delete(ctx context.Context, userId, postId int64) error
}

type PostService interface {
	Create(ctx context.Context, userId int64, createPost *domain.CreatePostDTO) (*domain.Post, error)
	List(ctx context.Context, limit int, offset int) ([]domain.Post, error)
	GetByID(ctx context.Context, postId int64) (*domain.Post, error)
	Update(ctx context.Context, userId, postId int64, post *domain.UpdatePostDTO) (*domain.Post, error)
	Delete(ctx context.Context, userId, postId int64) error
}
