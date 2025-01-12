package mocks

import (
	"context"

	"github.com/floroz/go-social/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockedUserRepository struct {
	mock.Mock
}

func (m *MockedUserRepository) Create(ctx context.Context, createUser *domain.CreateUserDTO) (*domain.User, error) {
	args := m.Called(ctx, createUser)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockedUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockedUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockedUserRepository) Update(ctx context.Context, userId int64, updateUser *domain.UpdateUserDTO) (*domain.User, error) {
	args := m.Called(ctx, userId, updateUser)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockedUserRepository) GetByID(ctx context.Context, userId int64) (*domain.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockedUserRepository) Delete(ctx context.Context, userId int64) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockedUserRepository) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.User), args.Error(1)
}
