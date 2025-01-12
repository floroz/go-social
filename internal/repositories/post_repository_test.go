package repositories_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/mocks"
	"github.com/floroz/go-social/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestPostRepositoryImpl_Create_Success(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	createPostDTO := &domain.CreatePostDTO{
		EditablePostFields: domain.EditablePostFields{
			Content: "Post Content",
		},
	}

	expectedPost := &domain.Post{
		ID:        1,
		UserID:    1,
		Content:   "Post Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`INSERT INTO posts`).
		WithArgs(expectedPost.UserID, createPostDTO.Content).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "content", "created_at", "updated_at"}).
			AddRow(expectedPost.ID, expectedPost.UserID, expectedPost.Content, expectedPost.CreatedAt, expectedPost.UpdatedAt))

	// Act
	post, err := repo.Create(context.Background(), expectedPost.UserID, createPostDTO)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, expectedPost, post)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepositoryImpl_Create_Error(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	createPostDTO := &domain.CreatePostDTO{
		EditablePostFields: domain.EditablePostFields{
			Content: "Post Content",
		},
	}
	mock.ExpectQuery("INSERT INTO posts").
		WithArgs(int64(1), createPostDTO.Content).
		WillReturnError(errors.New("some error"))

	// Act
	post, err := repo.Create(context.Background(), int64(1), createPostDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, post)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepositoryImpl_GetByID_Success(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	const postId int64 = 1
	expectedPost := &domain.Post{
		ID:        postId,
		UserID:    1,
		Content:   "Post Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`SELECT id, user_id, content, created_at, updated_at FROM posts WHERE id = \$1`).
		WithArgs(postId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "content", "created_at", "updated_at"}).
			AddRow(expectedPost.ID, expectedPost.UserID, expectedPost.Content, expectedPost.CreatedAt, expectedPost.UpdatedAt))

	// Act
	post, err := repo.GetByID(context.Background(), postId)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, expectedPost, post)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepositoryImpl_GetByID_Error(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	const postId int64 = 1

	mock.ExpectQuery(`SELECT id, user_id, content, created_at, updated_at FROM posts WHERE id = \$1`).
		WithArgs(postId).
		WillReturnError(errors.New("some error"))

	// Act
	post, err := repo.GetByID(context.Background(), postId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, post)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepositoryImpl_Delete_Success(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	const postId int64 = 1
	const userId int64 = 1

	mock.ExpectExec(`DELETE FROM posts WHERE id = \$1 AND user_id = \$2`).
		WithArgs(postId, userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), userId, postId)

	// Assert
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepositoryImpl_Delete_Error(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	const postId int64 = 1
	const userId int64 = 1

	mock.ExpectExec(`DELETE FROM posts WHERE id = \$1 AND user_id = \$2`).
		WithArgs(postId, userId).
		WillReturnError(errors.New("some error"))

	// Act
	err := repo.Delete(context.Background(), userId, postId)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepositoryImpl_List_Success(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	const limit, offset = 10, 0
	expectedPosts := []domain.Post{
		{ID: 1, UserID: 1, Content: "Content 1"},
		{ID: 2, UserID: 2, Content: "Content 2"},
	}

	post1 := expectedPosts[0]
	post2 := expectedPosts[1]

	mock.ExpectQuery(`SELECT id, user_id, content, created_at, updated_at FROM posts LIMIT \$1 OFFSET \$2`).
		WithArgs(limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "content", "created_at", "updated_at"}).
			AddRow(post1.ID, post1.UserID, post1.Content, post1.CreatedAt, post1.UpdatedAt).
			AddRow(post2.ID, post2.UserID, post2.Content, post2.CreatedAt, post2.UpdatedAt))

	// Act
	posts, err := repo.List(context.Background(), limit, offset)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, posts)
	assert.Equal(t, expectedPosts, posts)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepositoryImpl_List_Error(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	const limit, offset = 10, 0

	mock.ExpectQuery(`SELECT id, user_id, content, created_at, updated_at FROM posts LIMIT \$1 OFFSET \$2`).
		WithArgs(limit, offset).
		WillReturnError(errors.New("some error"))

	// Act
	posts, err := repo.List(context.Background(), limit, offset)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.NoError(t, mock.ExpectationsWereMet())
}
