package services

import (
	"context"

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
		return nil, domain.NewBadRequestError(err.Error())
	}

	// check existing email
	if existingUser, err := s.userRepo.GetByEmail(ctx, createUser.Email); existingUser != nil {
		// obfuscate error message to avoid leaking user information
		return nil, domain.NewBadRequestError("invalid body request")
	} else if err != nil && err != domain.ErrNotFound {
		log.Error().Err(err).Msg("failed to get user by email")
		return nil, domain.NewInternalServerError("something went wrong")
	}
	// check existing username
	if existingUser, err := s.userRepo.GetByUsername(ctx, createUser.Username); existingUser != nil {
		// obfuscate error message to avoid leaking user information
		return nil, domain.NewBadRequestError("invalid body request")
	} else if err != nil && err != domain.ErrNotFound {
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

func (s *userService) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)

	if err != nil && err == domain.ErrNotFound {
		log.Error().Err(err).Msg("failed to get user")
		return nil, domain.NewInternalServerError("failed to get user")
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, updateUser *domain.UpdateUserDTO) (*domain.User, error) {
	err := validation.Validate.Struct(updateUser)
	// TODO: check that current user is the owner of the user to be updated - after implementing authn and authz
	updatedUser, err := s.userRepo.Update(ctx, updateUser)

	if err != nil {
		log.Error().Err(err).Msg("failed to update user")
		return nil, domain.NewInternalServerError("failed to update user")
	}

	return updatedUser, err
}

func (s *userService) Delete(ctx context.Context, id int) error {
	// TODO: check that current user is the owner of the user to be deleted - after implementing authn and authz

	err := s.userRepo.Delete(ctx, id)

	if err != nil {
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
