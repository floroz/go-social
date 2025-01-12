package repositories_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, mock, cleanup
}

func TestUserRepositoryImpl_Create_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Username:  "johndoe",
		},
		Password: "hashedpassword",
	}

	expectedUser := &domain.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Username:  "johndoe",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(createUserDTO.FirstName, createUserDTO.LastName, createUserDTO.Email, createUserDTO.Username, createUserDTO.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "username", "password", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.Username, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt))

	// Act
	user, err := repo.Create(context.Background(), createUserDTO)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.FirstName, user.FirstName)
	assert.Equal(t, expectedUser.LastName, user.LastName)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Password, user.Password)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryImpl_Create_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	createUserDTO := &domain.CreateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Username:  "johndoe",
		},
		Password: "hashedpassword",
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(createUserDTO.FirstName, createUserDTO.LastName, createUserDTO.Email, createUserDTO.Username, createUserDTO.Password).
		WillReturnError(errors.New("some error"))

	// Act
	user, err := repo.Create(context.Background(), createUserDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryImpl_GetByID_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	const userId int64 = 1
	expectedUser := &domain.User{
		ID:        userId,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Username:  "johndoe",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, username, password, created_at, updated_at FROM users WHERE id = \$1`).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "username", "password", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.Username, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt))

	// Act
	user, err := repo.GetByID(context.Background(), userId)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryImpl_GetByID_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	const userId int64 = 1

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, username, password, created_at, updated_at FROM users WHERE id = \$1`).
		WithArgs(userId).
		WillReturnError(errors.New("some error"))

	// Act
	user, err := repo.GetByID(context.Background(), userId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryImpl_Update_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	const userId int64 = 1
	updateUserDTO := &domain.UpdateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Username:  "johndoe",
		},
	}

	expectedUser := &domain.User{
		ID:        userId,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Username:  "johndoe",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`UPDATE users SET first_name = \$1, last_name = \$2, email = \$3, username = \$4 WHERE id = \$5 RETURNING id, first_name, last_name, email, username, password, created_at, updated_at`).
		WithArgs(updateUserDTO.FirstName, updateUserDTO.LastName, updateUserDTO.Email, updateUserDTO.Username, userId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "username", "password", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.Username, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt))

	// Act
	user, err := repo.Update(context.Background(), userId, updateUserDTO)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryImpl_Update_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	const userId int64 = 1
	updateUserDTO := &domain.UpdateUserDTO{
		EditableUserField: domain.EditableUserField{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Username:  "johndoe",
		},
	}

	mock.ExpectQuery(`UPDATE users SET first_name = \$1, last_name = \$2, email = \$3, username = \$4 WHERE id = \$5 RETURNING id, first_name, last_name, email, username, password, created_at, updated_at`).
		WithArgs(updateUserDTO.FirstName, updateUserDTO.LastName, updateUserDTO.Email, updateUserDTO.Username, userId).
		WillReturnError(errors.New("some error"))

	// Act
	user, err := repo.Update(context.Background(), userId, updateUserDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryImpl_Delete_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	const userId int64 = 1

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), userId)

	// Assert
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryImpl_Delete_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	const userId int64 = 1

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs(userId).
		WillReturnError(errors.New("some error"))

	// Act
	err := repo.Delete(context.Background(), userId)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryImpl_List_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	const limit, offset = 10, 0
	expectedUsers := []domain.User{
		{ID: 1, FirstName: "Test1", LastName: "User1", Email: "test1@test.com", Username: "test1"},
		{ID: 2, FirstName: "Test2", LastName: "User2", Email: "test2@test.com", Username: "test2"},
	}

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, username, password, created_at, updated_at FROM users LIMIT \$1 OFFSET \$2`).
		WithArgs(limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "username", "password", "created_at", "updated_at"}).
			AddRow(expectedUsers[0].ID, expectedUsers[0].FirstName, expectedUsers[0].LastName, expectedUsers[0].Email, expectedUsers[0].Username, expectedUsers[0].Password, expectedUsers[0].CreatedAt, expectedUsers[0].UpdatedAt).
			AddRow(expectedUsers[1].ID, expectedUsers[1].FirstName, expectedUsers[1].LastName, expectedUsers[1].Email, expectedUsers[1].Username, expectedUsers[1].Password, expectedUsers[1].CreatedAt, expectedUsers[1].UpdatedAt))

	// Act
	users, err := repo.List(context.Background(), limit, offset)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, expectedUsers, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryImpl_List_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	const limit, offset = 10, 0

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, username, password, created_at, updated_at FROM users LIMIT \$1 OFFSET \$2`).
		WithArgs(limit, offset).
		WillReturnError(errors.New("some error"))

	// Act
	users, err := repo.List(context.Background(), limit, offset)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}
