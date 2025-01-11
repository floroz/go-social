package services

import (
	"context"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
	"github.com/floroz/go-social/internal/validation"
)

type commentsService struct {
	commentsRepo interfaces.CommentRepository
}

func NewCommentService(commentsRepo interfaces.CommentRepository) interfaces.CommentService {
	return &commentsService{commentsRepo: commentsRepo}
}

func (s *commentsService) Create(ctx context.Context, userId, postId int64, comment *domain.CreateCommentDTO) (*domain.Comment, error) {
	err := validation.Validate.Struct(comment)

	if err != nil {
		return nil, domain.NewBadRequestError(err.Error())
	}
	return s.commentsRepo.Create(ctx, userId, postId, comment)
}

func (s *commentsService) Delete(ctx context.Context, userId, commentId int64) error {
	comment, err := s.commentsRepo.GetByID(ctx, commentId)
	switch {
	case err != nil && err == domain.ErrNotFound:
		return domain.NewNotFoundError("comment not found")
	case err != nil:
		return domain.NewInternalServerError("failed to delete comment")
	case comment.UserID != userId:
		return domain.NewForbiddenError("not allowed to delete comment")
	}

	err = s.commentsRepo.Delete(ctx, userId, commentId)
	switch {
	case err != nil && err == domain.ErrNotFound:
		return domain.NewNotFoundError("comment not found")
	case err != nil:
		return domain.NewInternalServerError("failed to delete comment")
	default:
		return nil
	}

}

func (s *commentsService) GetByID(ctx context.Context, id int64) (*domain.Comment, error) {
	// todo validate authorization to get a comment

	return s.commentsRepo.GetByID(ctx, id)
}

func (s *commentsService) ListByPostID(ctx context.Context, postId int64, limit int, offset int) ([]domain.Comment, error) {
	// todo validate authorization to list comments for a post

	comments, err := s.commentsRepo.ListByPostID(ctx, postId, limit, offset)

	if err != nil && err == domain.ErrNotFound {
		return nil, domain.NewNotFoundError("post not found")
	}

	if err != nil {
		return nil, domain.NewInternalServerError("failed to list comments")
	}

	return comments, nil
}

func (s *commentsService) Update(ctx context.Context, userId, postId, commentId int64, comment *domain.UpdateCommentDTO) (*domain.Comment, error) {
	err := validation.Validate.Struct(comment)
	if err != nil {
		return nil, domain.NewBadRequestError(err.Error())
	}

	switch existing, err := s.commentsRepo.GetByID(ctx, commentId); {
	case err != nil && err == domain.ErrNotFound:
		return nil, domain.NewNotFoundError("comment not found")
	case err != nil:
		return nil, domain.NewInternalServerError("failed to delete comment")
	case existing.UserID != userId:
		return nil, domain.NewForbiddenError("not allowed to delete comment")
	}

	updatedComment, err := s.commentsRepo.Update(ctx, userId, postId, comment)

	if err != nil && err == domain.ErrNotFound {
		return nil, domain.NewNotFoundError("comment not found")
	}

	if err != nil {
		return nil, domain.NewInternalServerError("failed to update comment")
	}

	return updatedComment, nil
}
