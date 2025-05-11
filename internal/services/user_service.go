package services

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
	"github.com/floroz/go-social/internal/validation"
	"github.com/rs/zerolog/log"
)

type userService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) interfaces.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(ctx context.Context, createUser *domain.CreateUserDTO) (*domain.User, error) {
	err := validation.Validate.Struct(createUser)
	if err != nil {
		return nil, err
	}

	// check existing email
	if existingUser, err := s.userRepo.GetByEmail(ctx, createUser.Email); existingUser != nil {
		return nil, domain.ErrDuplicateEmailOrUsername
	} else if err != nil && !errors.Is(err, domain.ErrNotFound) {
		log.Error().Err(err).Msg("failed to get user by email")
		return nil, domain.NewInternalServerError("something went wrong")
	}
	// check existing username
	if existingUser, err := s.userRepo.GetByUsername(ctx, createUser.Username); existingUser != nil {
		return nil, domain.ErrDuplicateEmailOrUsername
	} else if err != nil && !errors.Is(err, domain.ErrNotFound) {
		log.Error().Err(err).Msg("failed to get user by username")
		return nil, domain.NewInternalServerError("something went wrong")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return nil, domain.NewInternalServerError("failed to hash password")
	}
	createUser.Password = string(hashedPassword)

	createdUser, err := s.userRepo.Create(ctx, createUser)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return nil, domain.NewInternalServerError("failed to create user")
	}

	return createdUser, nil
}

func (s *userService) GetByID(ctx context.Context, userId int64) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil && err == domain.ErrNotFound {
		log.Error().Err(err).Msg("failed to get user")
		return nil, domain.NewInternalServerError("failed to get user")
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, userId int64, updateUser *domain.UpdateUserDTO) (*domain.User, error) {
	if err := validation.Validate.Struct(updateUser); err != nil {
		return nil, domain.NewBadRequestError(err.Error())
	}

	// check a user exists first (original logic before partial update attempt)
	if _, err := s.userRepo.GetByID(ctx, userId); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.NewNotFoundError("user not found")
		}
		log.Error().Err(err).Msg("failed to get user by id before update")
		return nil, domain.NewInternalServerError("failed to get user by id")
	}

	// TODO: Add checks for existing email/username if they are being updated,
	// similar to the Create method, but only if the DTO values differ from the existing ones.
	// This requires fetching the user data again or modifying the GetByID check above.
	// For now, we rely on database constraints for uniqueness.

	// Call repository update directly with the validated DTO
	updatedUser, err := s.userRepo.Update(ctx, userId, updateUser)
	if err != nil {
		// Handle potential ErrNotFound from the repo update itself
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.NewNotFoundError("user not found during update")
		}
		log.Error().Err(err).Msg("failed to update user")
		return nil, domain.NewInternalServerError("failed to update user")
	}

	return updatedUser, err
}

func (s *userService) Delete(ctx context.Context, userId int64) error {
	err := s.userRepo.Delete(ctx, userId)

	if err == domain.ErrNotFound {
		return domain.NewNotFoundError("user not found")
	} else if err != nil {
		log.Error().Err(err).Msg("failed to delete user")
		return domain.NewInternalServerError("failed to delete user")
	}

	return nil
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to list users")
		return nil, domain.NewInternalServerError("failed to list users")
	}

	return users, nil
}

func (s *userService) UpdateLastLogin(ctx context.Context, userId int64) error {
	err := s.userRepo.UpdateLastLogin(ctx, userId)
	if err != nil {
		log.Error().Err(err).Msg("failed to update last login")
		return domain.NewInternalServerError("failed to update last login")
	}

	return nil
}
