package mocks

import (
	"context"

	"github.com/floroz/go-social/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockedPostRepository struct {
	mock.Mock
}

func (m *MockedPostRepository) Create(ctx context.Context, userId int64, post *domain.CreatePostDTO) (*domain.Post, error) {
	args := m.Called(ctx, userId, post)
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockedPostRepository) List(ctx context.Context, limit int, offset int) ([]*domain.Post, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*domain.Post), args.Error(1)
}

func (m *MockedPostRepository) GetByID(ctx context.Context, postId int64) (*domain.Post, error) {
	args := m.Called(ctx, postId)
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockedPostRepository) Update(ctx context.Context, userId, postId int64, post *domain.UpdatePostDTO) (*domain.Post, error) {
	args := m.Called(ctx, userId, postId, post)
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockedPostRepository) Delete(ctx context.Context, userId, postId int64) error {
	args := m.Called(ctx, userId, postId)
	return args.Error(0)
}
