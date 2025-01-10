package interfaces

import (
	"context"

	"github.com/floroz/go-social/internal/domain"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *domain.CreateCommentDTO) (*domain.Comment, error)
	Delete(ctx context.Context, id int) error
	DeleteByPostID(ctx context.Context, postId int) error
	GetByID(ctx context.Context, id int) (*domain.Comment, error)
	ListByPostID(ctx context.Context, postId int, limit int, offset int) ([]*domain.Comment, error)
	Update(ctx context.Context, comment *domain.UpdateCommentDTO) (*domain.Comment, error)
}

type CommentService interface {
	Create(ctx context.Context, comment *domain.CreateCommentDTO) (*domain.Comment, error)
	Delete(ctx context.Context, id int) error
	DeleteByPostID(ctx context.Context, postId int) error
	GetByID(ctx context.Context, id int) (*domain.Comment, error)
	ListByPostID(ctx context.Context, postId int, limit int, offset int) ([]*domain.Comment, error)
	Update(ctx context.Context, comment *domain.UpdateCommentDTO) (*domain.Comment, error)
}
