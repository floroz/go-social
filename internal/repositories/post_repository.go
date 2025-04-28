package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
)

type PostRepositoryImpl struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) interfaces.PostRepository {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) Create(ctx context.Context, userId int64, createPost *domain.CreatePostDTO) (*domain.Post, error) {
	query := `
		INSERT INTO posts (user_id, content)
		VALUES ($1, $2)
		RETURNING id, user_id, content, created_at, updated_at
		`

	newPost := domain.Post{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		userId,
		createPost.Content,
	).Scan(
		&newPost.ID,
		&newPost.UserID,
		&newPost.Content,
		&newPost.CreatedAt,
		&newPost.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &newPost, nil
}

func (r *PostRepositoryImpl) List(ctx context.Context, limit int, offset int) ([]domain.Post, error) {
	query := `
		SELECT id, user_id, content, created_at, updated_at
		FROM posts
		LIMIT $1 OFFSET $2
		`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := make([]domain.Post, 0)

	for rows.Next() {
		post := domain.Post{}

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepositoryImpl) GetByID(ctx context.Context, postId int64) (*domain.Post, error) {
	query := `
		SELECT id, user_id, content, created_at, updated_at
		FROM posts
		WHERE id = $1
		`

	post := domain.Post{}

	err := r.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &post, nil
}

func (r *PostRepositoryImpl) Update(ctx context.Context, userId, postId int64, post *domain.UpdatePostDTO) (*domain.Post, error) {
	query := `
		UPDATE posts
		SET content = $1
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, content, created_at, updated_at
		`

	updatedPost := domain.Post{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		postId,
		userId,
	).Scan(
		&updatedPost.ID,
		&updatedPost.UserID,
		&updatedPost.Content,
		&updatedPost.CreatedAt,
		&updatedPost.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &updatedPost, nil
}

func (r *PostRepositoryImpl) Delete(ctx context.Context, userId, postId int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1 AND user_id = $2
		`

	_, err := r.db.ExecContext(ctx, query, postId, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}
