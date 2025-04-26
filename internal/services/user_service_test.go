package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/mocks"
	"github.com/floroz/go-social/internal/services"
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
	mockUserRepo.On("GetByEmail", context.Background(), createUserDTO.Email).Return(&domain.User{}, nil)
	userService := services.NewUserService(mockUserRepo)

	// Act
	_, err := userService.Create(context.Background(), createUserDTO)

	// Assert
	assert.NotNil(t, err)
	mockUserRepo.AssertCalled(t, "GetByEmail", mock.Anything, createUserDTO.Email)
	mockUserRepo.AssertNotCalled(t, "GetByUsername")
	mockUserRepo.AssertNotCalled(t, "Create")
	mockUserRepo.AssertExpectations(t)
	assert.IsType(t, &domain.BadRequestError{}, err)
	assert.ErrorContains(t, err, "invalid body request")
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
	mockUserRepo.On("GetByEmail", mock.Anything, createUserDTO.Email).Return(existingUserWithEmail, domain.ErrNotFound)

	mockUserRepo.On("GetByUsername", mock.Anything, createUserDTO.Username).Return(&domain.User{}, nil)

	// Act
	_, err := userService.Create(context.Background(), createUserDTO)

	// Assert
	assert.NotNil(t, err)
	mockUserRepo.AssertCalled(t, "GetByEmail", mock.Anything, createUserDTO.Email)
	mockUserRepo.AssertCalled(t, "GetByUsername", mock.Anything, createUserDTO.Username)
	mockUserRepo.AssertNotCalled(t, "Create")
	mockUserRepo.AssertExpectations(t)
	assert.IsType(t, &domain.BadRequestError{}, err)
	assert.ErrorContains(t, err, "invalid body request")
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
	mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nullptr, errors.New("something went wrong"))

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

		// Assert password is hashed
		assert.NotNil(t, user)
		assert.NotEmpty(t, user.Password)
		assert.NotEqual(t, originalPassword, user.Password)
		assert.Greater(t, len(user.Password), len(originalPassword))
	}).Return(&domain.User{
		FirstName: createUserDTO.FirstName,
		LastName:  createUserDTO.LastName,
		Email:     createUserDTO.Email,
		Username:  createUserDTO.Username,
		Password:  createUserDTO.Password,
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
	mockUserRepo.On("GetByID", mock.Anything, userId).Return(nullptr, domain.ErrNotFound)

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
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@test.com",
			Username:  "test",
		},
	}
	mockUserRepo := new(mocks.MockedUserRepository)
	userService := services.NewUserService(mockUserRepo)

	const userId int64 = 1

	mockUserRepo.On("GetByID", mock.Anything, userId).Return(&domain.User{}, nil)
	mockUserRepo.On("Update", mock.Anything, userId, mock.Anything).Return(&domain.User{
		FirstName: updateUserDTO.FirstName,
		LastName:  updateUserDTO.LastName,
	}, nil)

	// Act
	user, err := userService.Update(context.Background(), 1, updateUserDTO)

	// Assert
	assert.Nil(t, err)
	mockUserRepo.AssertCalled(t, "GetByID", mock.Anything, userId)
	mockUserRepo.AssertCalled(t, "Update", mock.Anything, userId, mock.Anything)
	mockUserRepo.AssertExpectations(t)
	assert.NotNil(t, user)
	assert.Equal(t, updateUserDTO.FirstName, user.FirstName)
	assert.Equal(t, updateUserDTO.LastName, user.LastName)
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
					Email:     "test@test",
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
					Email:     strings.Repeat("A", 51) + "@test.com",
					Username:  "test",
				},
				Password: "password",
			},
			expectedError: "CreateUserDTO.EditableUserField.Email",
		},
		{
			name: "InvalidEmail - missing",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
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
					Username:  "test^^&",
				},
				Password: "password",
			},
			expectedError: "CreateUserDTO.EditableUserField.Username",
		},
		{
			name: "InvalidUsername - missing",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
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
				Password: "pass",
			},
			expectedError: "",
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
			},
			expectedError: "",
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
				Password: strings.Repeat("A", 51),
			},
			expectedError: "",
		},
		{
			name: "InvalidFirstName",
			createUserDTO: &domain.CreateUserDTO{
				EditableUserField: domain.EditableUserField{
					FirstName: "A",
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
					FirstName: strings.Repeat("A", 51),
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
				EditableUserField: domain.EditableUserField{
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
					LastName:  "U",
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
					LastName:  strings.Repeat("U", 51),
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
				EditableUserField: domain.EditableUserField{
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

			_, err := userService.Create(context.Background(), tt.createUserDTO)

			assert.NotNil(t, err)
			assert.IsType(t, &domain.BadRequestError{}, err)
			mockUserRepo.AssertNotCalled(t, "Create")
			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			}
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
