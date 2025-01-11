package interfaces

import (
	"context"

	"github.com/floroz/go-social/internal/domain"
)

type CommentRepository interface {
	Create(ctx context.Context, userId, postId int64, comment *domain.CreateCommentDTO) (*domain.Comment, error)
	Delete(ctx context.Context, userId, commentId int64) error
	GetByID(ctx context.Context, id int64) (*domain.Comment, error)
	ListByPostID(ctx context.Context, postId int64, limit int, offset int) ([]domain.Comment, error)
	Update(ctx context.Context, userId, postId int64, comment *domain.UpdateCommentDTO) (*domain.Comment, error)
}

type CommentService interface {
	Create(ctx context.Context, userId, postId int64, comment *domain.CreateCommentDTO) (*domain.Comment, error)
	Delete(ctx context.Context, userId, commentId int64) error
	GetByID(ctx context.Context, id int64) (*domain.Comment, error)
	ListByPostID(ctx context.Context, postId int64, limit int, offset int) ([]domain.Comment, error)
	Update(ctx context.Context, userId, postId, commentId int64, comment *domain.UpdateCommentDTO) (*domain.Comment, error)
}
