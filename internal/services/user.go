package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, createUser *domain.CreateUserDTO) (*domain.User, error)
	// GetUserByID(ctx context.Context, id int) (*domain.User, error)
	// UpdateUser(ctx context.Context, updateUser *domain.UpdateUserDTO) (*domain.User, error)
	// DeleteUser(ctx context.Context, id int) error
	// ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error)
}

type userService struct {
	userRepo repositories.UserRepositoryImpl
}

func NewUserService(userRepo repositories.UserRepositoryImpl) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *domain.CreateUserDTO) (*domain.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Check if email already exists
	// existing, err := s.userRepo.GetByEmail(ctx, user.Email)
	// if err != nil && !errors.Is(err, repositories.ErrNotFound) {
	// 	return nil, err
	// }
	// if existing != nil {
	// 	return nil, repositories.ErrDuplicate
	// }

	return s.userRepo.Create(ctx, user)
}

// func (s *userService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
// 	return s.userRepo.GetByID(ctx, id)
// }

// func (s *userService) UpdateUser(ctx context.Context, user *domain.UpdateUserDTO) (*domain.UpdateUserDTO, error) {
// 	if user.Password != "" {
// 		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to hash password: %w", err)
// 		}
// 		user.Password = string(hashedPassword)
// 	}
// 	return s.userRepo.Update(ctx, user), nil
// }

// func (s *userService) DeleteUser(ctx context.Context, id int) error {
// 	return s.userRepo.Delete(ctx, id)
// }

// func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
// 	return s.userRepo.List(ctx, limit, offset)
// }
