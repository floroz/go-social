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

func (s *commentsService) Create(ctx context.Context, comment *domain.CreateCommentDTO) (*domain.Comment, error) {
	err := validation.Validate.Struct(comment)

	// validate authorization to create a comment

	if err != nil {
		return nil, domain.NewBadRequestError(err.Error())
	}
	return s.commentsRepo.Create(ctx, comment)
}

func (s *commentsService) Delete(ctx context.Context, id int) error {
	// todo validate authorization to delete a comment

	err := s.commentsRepo.Delete(ctx, id)

	if err != nil && err == domain.ErrNotFound {
		return domain.NewNotFoundError("comment not found")
	}

	if err != nil {
		return domain.NewInternalServerError("failed to delete comment")
	}

	return nil
}

func (s *commentsService) DeleteByPostID(ctx context.Context, postId int) error {
	// todo validate authorization to delete all comment for a post

	err := s.commentsRepo.DeleteByPostID(ctx, postId)

	if err != nil && err == domain.ErrNotFound {
		return domain.NewNotFoundError("post not found")
	}

	if err != nil {
		return domain.NewInternalServerError("failed to delete comments")
	}

	return nil
}

func (s *commentsService) GetByID(ctx context.Context, id int) (*domain.Comment, error) {
	// todo validate authorization to get a comment

	return s.commentsRepo.GetByID(ctx, id)
}

func (s *commentsService) ListByPostID(ctx context.Context, postId int, limit int, offset int) ([]*domain.Comment, error) {
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

func (s *commentsService) Update(ctx context.Context, comment *domain.UpdateCommentDTO) (*domain.Comment, error) {
	// todo validate authorization to update a comment

	err := validation.Validate.Struct(comment)

	if err != nil {
		return nil, domain.NewBadRequestError(err.Error())
	}

	updatedComment, err := s.commentsRepo.Update(ctx, comment)

	if err != nil && err == domain.ErrNotFound {
		return nil, domain.NewNotFoundError("comment not found")
	}

	if err != nil {
		return nil, domain.NewInternalServerError("failed to update comment")
	}

	return updatedComment, nil
}
