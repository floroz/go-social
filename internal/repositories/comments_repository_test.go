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

func TestCommentRepositoryImpl_Create_Success(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewCommentRepository(db)

	createCommentDTO := &domain.CreateCommentDTO{
		EditableCommentFields: domain.EditableCommentFields{
			Content: "Updated comment",
		},
	}

	expectedComment := &domain.Comment{
		ID:        1,
		UserID:    1,
		PostID:    1,
		Content:   "This is a comment",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`INSERT INTO comments`).
		WithArgs(expectedComment.UserID, expectedComment.PostID, createCommentDTO.Content).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "content", "created_at", "updated_at"}).
			AddRow(expectedComment.ID, expectedComment.UserID, expectedComment.PostID, expectedComment.Content, expectedComment.CreatedAt, expectedComment.UpdatedAt))

	// Act
	comment, err := repo.Create(context.Background(), expectedComment.UserID, expectedComment.PostID, createCommentDTO)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, expectedComment, comment)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepositoryImpl_Create_Error(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewCommentRepository(db)

	createCommentDTO := &domain.CreateCommentDTO{
		EditableCommentFields: domain.EditableCommentFields{
			Content: "Updated comment",
		},
	}

	mock.ExpectQuery("INSERT INTO comments").
		WithArgs(int64(1), int64(1), createCommentDTO.Content).
		WillReturnError(errors.New("some error"))

	// Act
	comment, err := repo.Create(context.Background(), int64(1), int64(1), createCommentDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, comment)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepositoryImpl_GetByID_Success(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewCommentRepository(db)

	const commentId int64 = 1
	expectedComment := &domain.Comment{
		ID:        commentId,
		UserID:    1,
		PostID:    1,
		Content:   "This is a comment",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`SELECT id, user_id, post_id, content, created_at, updated_at FROM comments WHERE id = \$1`).
		WithArgs(commentId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "content", "created_at", "updated_at"}).
			AddRow(expectedComment.ID, expectedComment.UserID, expectedComment.PostID, expectedComment.Content, expectedComment.CreatedAt, expectedComment.UpdatedAt))

	// Act
	comment, err := repo.GetByID(context.Background(), commentId)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, expectedComment, comment)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepositoryImpl_GetByID_Error(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewCommentRepository(db)

	const commentId int64 = 1

	mock.ExpectQuery(`SELECT id, user_id, post_id, content, created_at, updated_at FROM comments WHERE id = \$1`).
		WithArgs(commentId).
		WillReturnError(errors.New("some error"))

	// Act
	comment, err := repo.GetByID(context.Background(), commentId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, comment)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepositoryImpl_Delete_Success(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewCommentRepository(db)

	const commentId int64 = 1
	const userId int64 = 1

	mock.ExpectExec(`DELETE FROM comments WHERE id = \$1 AND user_id = \$2`).
		WithArgs(commentId, userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), userId, commentId)

	// Assert
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepositoryImpl_Delete_Error(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewCommentRepository(db)

	const commentId int64 = 1
	const userId int64 = 1

	mock.ExpectExec(`DELETE FROM comments WHERE id = \$1 AND user_id = \$2`).
		WithArgs(commentId, userId).
		WillReturnError(errors.New("some error"))

	// Act
	err := repo.Delete(context.Background(), userId, commentId)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepositoryImpl_ListByPostID_Success(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewCommentRepository(db)

	const postId int64 = 1
	const limit, offset = 10, 0
	expectedComments := []domain.Comment{
		{ID: 1, UserID: 1, PostID: postId, Content: "Comment 1"},
		{ID: 2, UserID: 2, PostID: postId, Content: "Comment 2"},
	}

	mock.ExpectQuery(`SELECT id, user_id, post_id, content, created_at, updated_at FROM comments WHERE post_id = \$1 ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
		WithArgs(postId, limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "content", "created_at", "updated_at"}).
			AddRow(expectedComments[0].ID, expectedComments[0].UserID, expectedComments[0].PostID, expectedComments[0].Content, expectedComments[0].CreatedAt, expectedComments[0].UpdatedAt).
			AddRow(expectedComments[1].ID, expectedComments[1].UserID, expectedComments[1].PostID, expectedComments[1].Content, expectedComments[1].CreatedAt, expectedComments[1].UpdatedAt))

	// Act
	comments, err := repo.ListByPostID(context.Background(), postId, limit, offset)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, comments)
	assert.Equal(t, expectedComments, comments)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepositoryImpl_ListByPostID_Error(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewCommentRepository(db)

	const postId int64 = 1
	const limit, offset = 10, 0

	mock.ExpectQuery(`SELECT id, user_id, post_id, content, created_at, updated_at FROM comments WHERE post_id = \$1 ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
		WithArgs(postId, limit, offset).
		WillReturnError(errors.New("some error"))

	// Act
	comments, err := repo.ListByPostID(context.Background(), postId, limit, offset)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, comments)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepositoryImpl_Update_Success(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewCommentRepository(db)

	const commentId int64 = 1
	const userId int64 = 1
	updateCommentDTO := &domain.UpdateCommentDTO{
		ID: commentId,
		EditableCommentFields: domain.EditableCommentFields{
			Content: "Updated comment",
		},
	}

	expectedComment := &domain.Comment{
		ID:        commentId,
		UserID:    userId,
		PostID:    1,
		Content:   "Updated comment",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`UPDATE comments SET content = \$1 WHERE id = \$2 AND user_id = \$3 RETURNING id, user_id, post_id, content, created_at, updated_at`).
		WithArgs(updateCommentDTO.Content, updateCommentDTO.ID, userId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "content", "created_at", "updated_at"}).
			AddRow(expectedComment.ID, expectedComment.UserID, expectedComment.PostID, expectedComment.Content, expectedComment.CreatedAt, expectedComment.UpdatedAt))

	// Act
	comment, err := repo.Update(context.Background(), userId, expectedComment.PostID, updateCommentDTO)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, expectedComment, comment)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepositoryImpl_Update_Error(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewCommentRepository(db)

	const commentId int64 = 1
	const userId int64 = 1
	updateCommentDTO := &domain.UpdateCommentDTO{
		ID: commentId,
		EditableCommentFields: domain.EditableCommentFields{
			Content: "Updated comment",
		},
	}

	mock.ExpectQuery(`UPDATE comments SET content = \$1 WHERE id = \$2 AND user_id = \$3 RETURNING id, user_id, post_id, content, created_at, updated_at`).
		WithArgs(updateCommentDTO.Content, updateCommentDTO.ID, userId).
		WillReturnError(errors.New("some error"))

	// Act
	comment, err := repo.Update(context.Background(), userId, int64(1), updateCommentDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, comment)
	assert.NoError(t, mock.ExpectationsWereMet())
}
