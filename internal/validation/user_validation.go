package validation

import (
	"regexp"

	"github.com/floroz/go-social/internal/domain"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func validateEmail(email string) error {
	switch {
	case email == "":
		return domain.NewValidationError("email", "email is required")
	case !emailRegex.MatchString(email):
		return domain.NewValidationError("email", "invalid email format")
	}

	return nil
}

func validatePassword(password string) error {
	switch {
	case password == "":
		return domain.NewValidationError("password", "password is required")
	case len(password) < 8:
		return domain.NewValidationError("password", "password must be at least 8 characters")
	}

	return nil
}

func validateUsername(username string) error {
	switch {
	case username == "":
		return domain.NewValidationError("username", "username is required")
	case len(username) < 3:
		return domain.NewValidationError("username", "username must be at least 3 characters")
	case len(username) > 50:
		return domain.NewValidationError("username", "username must be at most 50 characters")
	}

	return nil
}

func validateFirstName(firstName string) error {
	switch {
	case firstName == "":
		return domain.NewValidationError("first_name", "first name is required")
	case len(firstName) < 3:
		return domain.NewValidationError("first_name", "first name must be at least 3 characters")
	case len(firstName) > 50:
		return domain.NewValidationError("first_name", "first name must be at most 50 characters")
	}

	return nil
}

func validateLastName(lastName string) error {
	switch {
	case lastName == "":
		return domain.NewValidationError("last_name", "last name is required")
	case len(lastName) < 3:
		return domain.NewValidationError("last_name", "last name must be at least 3 characters")
	case len(lastName) > 50:
		return domain.NewValidationError("last_name", "last name must be at most 50 characters")
	}

	return nil
}

func ValidateCreateUserDTO(user *domain.CreateUserDTO) error {

	if err := validateEmail(user.Email); err != nil {
		return err
	}

	if err := validatePassword(user.Password); err != nil {
		return err
	}

	if err := validateUsername(user.Username); err != nil {
		return err
	}

	if err := validateFirstName(user.FirstName); err != nil {
		return err
	}

	if err := validateLastName(user.LastName); err != nil {
		return err
	}

	return nil
}

func ValidateUpdateUserDTO(user *domain.UpdateUserDTO) error {
	if err := validateEmail(user.Email); err != nil {
		return err
	}

	if err := validateUsername(user.Username); err != nil {
		return err
	}

	if err := validateFirstName(user.FirstName); err != nil {
		return err
	}

	if err := validateLastName(user.LastName); err != nil {
		return err
	}

	return nil
}
