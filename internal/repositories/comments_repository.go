package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
)

type CommentRepositoryImpl struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) interfaces.CommentRepository {
	return &CommentRepositoryImpl{db: db}
}

func (r *CommentRepositoryImpl) Create(ctx context.Context, comment *domain.CreateCommentDTO) (*domain.Comment, error) {
	query := `
		INSERT INTO comments (user_id, post_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, post_id, content, created_at, updated_at
		`

	newComment := domain.Comment{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		comment.UserID,
		comment.PostID,
		comment.Content,
	).Scan(
		&newComment.ID,
		&newComment.UserID,
		&newComment.PostID,
		&newComment.Content,
		&newComment.CreatedAt,
		&newComment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &newComment, nil
}

func (r *CommentRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM comments
		WHERE id = $1
		`

	_, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}

		return err
	}

	return nil
}

func (r *CommentRepositoryImpl) DeleteByPostID(ctx context.Context, postId int) error {
	query := `
		DELETE FROM comments
		WHERE post_id = $1
		`

	_, err := r.db.ExecContext(ctx, query, postId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}

		return err
	}

	return nil
}

func (r *CommentRepositoryImpl) GetByID(ctx context.Context, id int) (*domain.Comment, error) {
	query := `
		SELECT id, user_id, post_id, content, created_at, updated_at
		FROM comments
		WHERE id = $1
		`

	comment := domain.Comment{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.UserID,
		&comment.PostID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepositoryImpl) ListByPostID(ctx context.Context, postId int, limit int, offset int) ([]*domain.Comment, error) {
	query := `
		SELECT id, user_id, post_id, content, created_at, updated_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at DESC
		LIMIT $2
		OFFSET $3
		`

	rows, err := r.db.QueryContext(ctx, query, postId, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := make([]*domain.Comment, 0)

	for rows.Next() {
		comment := domain.Comment{}

		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.PostID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}
