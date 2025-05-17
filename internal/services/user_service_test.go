package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/mocks"
	"github.com/floroz/go-social/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_ExistingEmail(t *testing.T) {
	// Arrange
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@test.com",
			Username:  "test",
		},
		Password: "password",
	}
	mockUserRepo := new(mocks.MockedUserRepository)
	mockUserRepo.On("GetByEmail", context.Background(), createUserDTO.Email).Return(&domain.User{}, nil) // Simulate existing email
	userService := services.NewUserService(mockUserRepo)

	// Act
	_, err := userService.Create(context.Background(), createUserDTO)

	// Assert
	assert.NotNil(t, err)
	mockUserRepo.AssertCalled(t, "GetByEmail", mock.Anything, createUserDTO.Email)
	mockUserRepo.AssertNotCalled(t, "GetByUsername")
	mockUserRepo.AssertNotCalled(t, "Create")
	mockUserRepo.AssertExpectations(t)
	assert.True(t, errors.Is(err, domain.ErrDuplicateEmailOrUsername), "Expected ErrDuplicateEmailOrUsername")
}

func TestCreateUser_ExistingUsername(t *testing.T) {
	// Arrange
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@test.com",
			Username:  "test",
		},
		Password: "password",
	}
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	var existingUserWithEmail *domain.User
	mockUserRepo.On("GetByEmail", mock.Anything, createUserDTO.Email).Return(existingUserWithEmail, domain.ErrNotFound) // Simulate email not found
	mockUserRepo.On("GetByUsername", mock.Anything, createUserDTO.Username).Return(&domain.User{}, nil)                 // Simulate username found

	// Act
	_, err := userService.Create(context.Background(), createUserDTO)

	// Assert
	assert.NotNil(t, err)
	mockUserRepo.AssertCalled(t, "GetByEmail", mock.Anything, createUserDTO.Email)
	mockUserRepo.AssertCalled(t, "GetByUsername", mock.Anything, createUserDTO.Username)
	mockUserRepo.AssertNotCalled(t, "Create")
	mockUserRepo.AssertExpectations(t)
	assert.True(t, errors.Is(err, domain.ErrDuplicateEmailOrUsername), "Expected ErrDuplicateEmailOrUsername")
}

func TestCreateUser_ErrorCreating(t *testing.T) {
	// Arrange
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@test.com",
			Username:  "test",
		},
		Password: "password",
	}
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	var nullptr *domain.User

	mockUserRepo.On("GetByEmail", mock.Anything, createUserDTO.Email).Return(nullptr, domain.ErrNotFound)
	mockUserRepo.On("GetByUsername", mock.Anything, createUserDTO.Username).Return(nullptr, domain.ErrNotFound)
	mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nullptr, errors.New("something went wrong")) // Simulate DB error on Create

	// Act
	_, err := userService.Create(context.Background(), createUserDTO)

	// Assert
	assert.NotNil(t, err)
	mockUserRepo.AssertCalled(t, "GetByEmail", mock.Anything, createUserDTO.Email)
	mockUserRepo.AssertCalled(t, "GetByUsername", mock.Anything, createUserDTO.Username)
	mockUserRepo.AssertCalled(t, "Create", mock.Anything, mock.Anything)
	mockUserRepo.AssertExpectations(t)
	assert.IsType(t, &domain.InternalServerError{}, err)
	assert.ErrorContains(t, err, "failed to create user")
}

func TestCreateUser_Success(t *testing.T) {
	const originalPassword = "password"
	// Arrange
	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@test.com",
			Username:  "test",
		},
		Password: originalPassword,
	}
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	var existingUser *domain.User
	mockUserRepo.On("GetByEmail", mock.Anything, createUserDTO.Email).Return(existingUser, domain.ErrNotFound)
	mockUserRepo.On("GetByUsername", mock.Anything, createUserDTO.Username).Return(existingUser, domain.ErrNotFound)

	mockUserRepo.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		user := args.Get(1).(*domain.CreateUserDTO)
		assert.NotNil(t, user)
		assert.NotEmpty(t, user.Password)
		assert.NotEqual(t, originalPassword, user.Password)
		assert.Greater(t, len(user.Password), len(originalPassword))
	}).Return(&domain.User{
		FirstName: createUserDTO.FirstName,
		LastName:  createUserDTO.LastName,
		Email:     createUserDTO.Email,
		Username:  createUserDTO.Username,
		Password:  createUserDTO.Password, // This will be the hashed one in reality, but fine for mock
	}, nil)

	// Act
	user, err := userService.Create(context.Background(), createUserDTO)

	// Assert
	assert.Nil(t, err)
	mockUserRepo.AssertCalled(t, "GetByEmail", mock.Anything, createUserDTO.Email)
	mockUserRepo.AssertCalled(t, "GetByUsername", mock.Anything, createUserDTO.Username)
	mockUserRepo.AssertCalled(t, "Create", mock.Anything, mock.Anything)
	mockUserRepo.AssertExpectations(t)
	assert.NotNil(t, user)
	assert.NotEqual(t, createUserDTO.Password, user.Password)
	assert.Equal(t, createUserDTO.FirstName, user.FirstName)
	assert.Equal(t, createUserDTO.LastName, user.LastName)
	assert.Equal(t, createUserDTO.Email, user.Email)
	assert.Equal(t, createUserDTO.Username, user.Username)
}

func TestUpdateUser_NotExisting(t *testing.T) {
	// Arrange
	updateUserDTO := &domain.UpdateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@test.com",
			Username:  "test",
		},
	}
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	const userId int64 = 1

	var nullptr *domain.User
	mockUserRepo.On("GetByID", mock.Anything, userId).Return(nullptr, domain.ErrNotFound) // Simulate user not found

	// Act
	user, err := userService.Update(context.Background(), 1, updateUserDTO)

	// Assert
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.IsType(t, &domain.NotFoundError{}, err)
	mockUserRepo.AssertCalled(t, "GetByID", mock.Anything, userId)
	mockUserRepo.AssertNotCalled(t, "Update")
	mockUserRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	// Arrange
	updateUserDTO := &domain.UpdateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "UpdatedFirst",
			LastName:  "UpdatedLast",
			Email:     "update@test.com",
			Username:  "updateduser",
		},
	}
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	const userId int64 = 1
	existingUser := &domain.User{ // Simulate the user fetched by GetByID
		ID:        userId,
		FirstName: "OriginalFirst",
		LastName:  "OriginalLast",
		Email:     "original@test.com",
		Username:  "originaluser",
	}

	mockUserRepo.On("GetByID", mock.Anything, userId).Return(existingUser, nil) // Expect GetByID first
	// Corrected: No GetByEmail/GetByUsername checks in the reverted Update logic
	// Expect the Update call directly with the input DTO
	mockUserRepo.On("Update", mock.Anything, userId, updateUserDTO).Return(&domain.User{ // Return the updated user data
		ID:        userId,
		FirstName: updateUserDTO.FirstName,
		LastName:  updateUserDTO.LastName,
		Email:     updateUserDTO.Email,
		Username:  updateUserDTO.Username,
	}, nil)

	// Act
	user, err := userService.Update(context.Background(), 1, updateUserDTO)

	// Assert
	assert.Nil(t, err)
	mockUserRepo.AssertCalled(t, "GetByID", mock.Anything, userId)
	// Corrected: Assert GetByEmail/GetByUsername are NOT called
	mockUserRepo.AssertNotCalled(t, "GetByEmail")
	mockUserRepo.AssertNotCalled(t, "GetByUsername")
	mockUserRepo.AssertCalled(t, "Update", mock.Anything, userId, updateUserDTO) // Assert Update is called with correct DTO
	mockUserRepo.AssertExpectations(t)
	assert.NotNil(t, user)
	assert.Equal(t, updateUserDTO.FirstName, user.FirstName)
	assert.Equal(t, updateUserDTO.LastName, user.LastName)
	assert.Equal(t, updateUserDTO.Email, user.Email)
	assert.Equal(t, updateUserDTO.Username, user.Username)
}

func TestCreateUser_Validation(t *testing.T) {
	tests := []struct {
		name          string
		createUserDTO *domain.CreateUserDTO
		expectedError string
	}{
		{
			name: "InvalidEmail",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "Test",
					LastName:  "User",
					Email:     "test@test", // Invalid format
					Username:  "test",
				},
				Password: "password",
			},
			expectedError: "CreateUserDTO.EditableUserField.Email",
		},
		{
			name: "InvalidEmail - too long",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "Test",
					LastName:  "User",
					Email:     strings.Repeat("A", 51) + "@test.com", // Too long
					Username:  "test",
				},
				Password: "password",
			},
			expectedError: "CreateUserDTO.EditableUserField.Email",
		},
		{
			name: "InvalidEmail - missing",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{ // Email is missing (zero value "")
					FirstName: "Test",
					LastName:  "User",
					Username:  "test",
				},
				Password: "password",
			},
			expectedError: "CreateUserDTO.EditableUserField.Email",
		},
		{
			name: "InvalidUsername",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "Test",
					LastName:  "User",
					Email:     "test@test.com",
					Username:  "test^^&", // Invalid chars
				},
				Password: "password",
			},
			expectedError: "CreateUserDTO.EditableUserField.Username",
		},
		{
			name: "InvalidUsername - missing",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{ // Username is missing
					FirstName: "Test",
					LastName:  "User",
					Email:     "test@test.com",
				},
				Password: "password",
			},
			expectedError: "CreateUserDTO.EditableUserField.Username",
		},
		{
			name: "InvalidPassword",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "Test",
					LastName:  "User",
					Email:     "test@test.com",
					Username:  "test",
				},
				Password: "pass", // Too short
			},
			expectedError: "CreateUserDTO.Password", // Expect error on Password field now
		},
		{
			name: "InvalidPassword - missing",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "Test",
					LastName:  "User",
					Email:     "test@test.com",
					Username:  "test",
				},
				// Password missing
			},
			expectedError: "CreateUserDTO.Password",
		},
		{
			name: "InvalidPassword - too long",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "Test",
					LastName:  "User",
					Email:     "test@test.com",
					Username:  "test",
				},
				Password: strings.Repeat("A", 51), // Too long
			},
			expectedError: "CreateUserDTO.Password",
		},
		{
			name: "InvalidFirstName",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "A", // Too short
					LastName:  "User",
					Email:     "test@test.com",
					Username:  "test",
				},
				Password: "password123",
			},
			expectedError: "CreateUserDTO.EditableUserField.FirstName",
		},
		{
			name: "InvalidFirstName - too long",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: strings.Repeat("A", 51), // Too long
					LastName:  "User",
					Email:     "test@test.com",
					Username:  "test",
				},
				Password: "password123",
			},
			expectedError: "CreateUserDTO.EditableUserField.FirstName",
		},
		{
			name: "InvalidFirstName - missing",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{ // FirstName missing
					LastName: "User",
					Email:    "test@test.com",
					Username: "test",
				},
				Password: "password123",
			},
			expectedError: "CreateUserDTO.EditableUserField.FirstName",
		},
		{
			name: "InvalidLastName",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "Antony",
					LastName:  "U", // Too short
					Email:     "test@test.com",
					Username:  "test",
				},
				Password: "password123",
			},
			expectedError: "CreateUserDTO.EditableUserField.LastName",
		},
		{
			name: "InvalidLastName - too long",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "Antony",
					LastName:  strings.Repeat("U", 51), // Too long
					Email:     "test@test.com",
					Username:  "test",
				},
				Password: "password123",
			},
			expectedError: "CreateUserDTO.EditableUserField.LastName",
		},
		{
			name: "InvalidLastName - missing",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{ // LastName missing
					FirstName: "Antony",
					Email:     "test@test.com",
					Username:  "test",
				},
				Password: "password123",
			},
			expectedError: "CreateUserDTO.EditableUserField.LastName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(mocks.MockedUserRepository)
			userService := services.NewUserService(mockUserRepo)

			// Validation happens before DB calls, so GetByEmail/GetByUsername should not be called.
			_, err := userService.Create(context.Background(), tt.createUserDTO)

			assert.NotNil(t, err)
			// Expect validator.ValidationErrors directly from the service
			var validationErrs validator.ValidationErrors
			assert.ErrorAs(t, err, &validationErrs, "Error should be validator.ValidationErrors")
			mockUserRepo.AssertNotCalled(t, "GetByEmail")
			mockUserRepo.AssertNotCalled(t, "GetByUsername")
			mockUserRepo.AssertNotCalled(t, "Create")
			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			}
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestGetUserByID_Success(t *testing.T) {
	// Arrange
	const userId int64 = 1
	expectedUser := &domain.User{
		ID:        userId,
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@test.com",
		Username:  "test",
	}
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	mockUserRepo.On("GetByID", mock.Anything, userId).Return(expectedUser, nil)

	// Act
	user, err := userService.GetByID(context.Background(), userId)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)
	mockUserRepo.AssertCalled(t, "GetByID", mock.Anything, userId)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	// Arrange
	const userId int64 = 1
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	var nullptr *domain.User
	mockUserRepo.On("GetByID", mock.Anything, userId).Return(nullptr, domain.ErrNotFound)

	// Act
	user, err := userService.GetByID(context.Background(), userId)

	// Assert
	assert.Nil(t, user)
	assert.NotNil(t, err)
	// Corrected: GetByID service layer returns InternalServerError if repo returns ErrNotFound
	assert.IsType(t, &domain.InternalServerError{}, err)
	mockUserRepo.AssertCalled(t, "GetByID", mock.Anything, userId)
	mockUserRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	// Arrange
	const userId int64 = 1
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	mockUserRepo.On("Delete", mock.Anything, userId).Return(nil)

	// Act
	err := userService.Delete(context.Background(), userId)

	// Assert
	assert.Nil(t, err)
	mockUserRepo.AssertCalled(t, "Delete", mock.Anything, userId)
	mockUserRepo.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
	// Arrange
	const userId int64 = 1
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	mockUserRepo.On("Delete", mock.Anything, userId).Return(domain.ErrNotFound)

	// Act
	err := userService.Delete(context.Background(), userId)

	// Assert
	assert.NotNil(t, err)
	assert.IsType(t, &domain.NotFoundError{}, err)
	mockUserRepo.AssertCalled(t, "Delete", mock.Anything, userId)
	mockUserRepo.AssertExpectations(t)
}

func TestListUsers_Success(t *testing.T) {
	// Arrange
	const limit, offset = 10, 0
	expectedUsers := []domain.User{
		{ID: 1, FirstName: "Test1", LastName: "User1", Email: "test1@test.com", Username: "test1"},
		{ID: 2, FirstName: "Test2", LastName: "User2", Email: "test2@test.com", Username: "test2"},
	}
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	mockUserRepo.On("List", mock.Anything, limit, offset).Return(expectedUsers, nil)

	// Act
	users, err := userService.List(context.Background(), limit, offset)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, expectedUsers, users)
	mockUserRepo.AssertCalled(t, "List", mock.Anything, limit, offset)
	mockUserRepo.AssertExpectations(t)
}

func TestListUsers_Error(t *testing.T) {
	// Arrange
	const limit, offset = 10, 0
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	var nullptr []domain.User
	mockUserRepo.On("List", mock.Anything, limit, offset).Return(nullptr, errors.New("something went wrong"))

	// Act
	users, err := userService.List(context.Background(), limit, offset)

	// Assert
	assert.Nil(t, users)
	assert.NotNil(t, err)
	assert.IsType(t, &domain.InternalServerError{}, err)
	mockUserRepo.AssertCalled(t, "List", mock.Anything, limit, offset)
	mockUserRepo.AssertExpectations(t)
}
