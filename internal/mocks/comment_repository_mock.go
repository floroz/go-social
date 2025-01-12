package mocks

import (
	"context"

	"github.com/floroz/go-social/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockedCommentRepository struct {
	mock.Mock
}

func (m *MockedCommentRepository) Create(ctx context.Context, userId, postId int64, comment *domain.CreateCommentDTO) (*domain.Comment, error) {
	args := m.Called(ctx, userId, postId, comment)
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockedCommentRepository) Delete(ctx context.Context, userId, commentId int64) error {
	args := m.Called(ctx, userId, commentId)
	return args.Error(0)
}

func (m *MockedCommentRepository) GetByID(ctx context.Context, id int64) (*domain.Comment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockedCommentRepository) ListByPostID(ctx context.Context, postId int64, limit int, offset int) ([]domain.Comment, error) {
	args := m.Called(ctx, postId, limit, offset)
	return args.Get(0).([]domain.Comment), args.Error(1)
}

func (m *MockedCommentRepository) Update(ctx context.Context, userId, postId int64, comment *domain.UpdateCommentDTO) (*domain.Comment, error) {
	args := m.Called(ctx, userId, postId, comment)
	return args.Get(0).(*domain.Comment), args.Error(1)
}
