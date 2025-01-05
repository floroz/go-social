package services

import (
	"context"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type userService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) interfaces.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *domain.CreateUserDTO) (*domain.User, error) {
	// validation
	switch {
	case user.FirstName == "":
		return nil, domain.NewValidationError("first_name", "first name is required")
	case user.LastName == "":
		return nil, domain.NewValidationError("last_name", "last name is required")
	case user.Email == "":
		return nil, domain.NewValidationError("email", "email is required")
	case !emailRegex.MatchString(user.Email):
		return nil, domain.NewValidationError("email", "invalid email format")
	case user.Password == "":
		return nil, domain.NewValidationError("password", "password is required")
	case len(user.Password) < 8:
		return nil, domain.NewValidationError("password", "password must be at least 8 characters")
	}

	// check that user with email does not exist
	existing, err := s.userRepo.GetByEmail(ctx, user.Email)

	if existing != nil {
		// obfuscate error message to avoid leaking user information
		return nil, domain.NewBadRequestError("invalid body request")
	}

	if err != nil && err != domain.ErrNotFound {
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
