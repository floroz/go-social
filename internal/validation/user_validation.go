package validation

import (
	"regexp"

	"github.com/floroz/go-social/internal/domain"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserDTO(user *domain.CreateUserDTO) error {
	switch {
	case user.FirstName == "":
		return domain.NewValidationError("first_name", "first name is required")
	case len(user.FirstName) < 3:
		return domain.NewValidationError("first_name", "first name must be at least 3 characters")
	case len(user.FirstName) > 50:
		return domain.NewValidationError("first_name", "first name must be at most 50 characters")
	case user.LastName == "":
		return domain.NewValidationError("last_name", "last name is required")
	case len(user.LastName) < 3:
		return domain.NewValidationError("last_name", "last name must be at least 3 characters")
	case len(user.LastName) > 50:
		return domain.NewValidationError("last_name", "last name must be at most 50 characters")
	case user.Email == "":
		return domain.NewValidationError("email", "email is required")
	case !emailRegex.MatchString(user.Email):
		return domain.NewValidationError("email", "invalid email format")
	case user.Password == "":
		return domain.NewValidationError("password", "password is required")
	case len(user.Password) < 8:
		return domain.NewValidationError("password", "password must be at least 8 characters")
	default:
		return nil
	}
}
