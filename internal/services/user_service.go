package services

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
	"github.com/floroz/go-social/internal/validation"
)

type userService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) interfaces.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(ctx context.Context, user *domain.CreateUserDTO) (*domain.User, error) {
	err := validation.ValidateUserDTO(user)
	if err != nil {
		return nil, err
	}

	// check existing email
	if existingUser, err := s.userRepo.GetByEmail(ctx, user.Email); existingUser != nil {
		// obfuscate error message to avoid leaking user information
		return nil, domain.NewBadRequestError("invalid body request")
	} else if err != nil && err != domain.ErrNotFound {
		return nil, domain.NewInternalServerError("something went wrong")
	}
	// check existing username
	if existingUser, err := s.userRepo.GetByUsername(ctx, user.Username); existingUser != nil {
		// obfuscate error message to avoid leaking user information
		return nil, domain.NewBadRequestError("invalid body request")
	} else if err != nil && err != domain.ErrNotFound {
		return nil, domain.NewInternalServerError("something went wrong")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.NewInternalServerError("failed to hash password")
	}

	// replace plain password with hashed password
	user.Password = string(hashedPassword)

	// Create user
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *userService) GetByID(ctx context.Context, id int) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, user *domain.UpdateUserDTO) (*domain.User, error) {
	updatedUser, err := s.userRepo.Update(ctx, user)
	return updatedUser, err
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return s.userRepo.List(ctx, limit, offset)
}
