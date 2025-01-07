package repositories

import (
	"context"
	"database/sql"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
)

type PostRepositoryImpl struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) interfaces.PostRepository {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) Create(ctx context.Context, post *domain.CreatePostDTO) (*domain.Post, error) {
	query := `
		INSERT INTO posts (user_id, content)
		VALUES ($1, $2)
		RETURNING id, user_id, content, created_at, updated_at
		`

	newPost := domain.Post{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		post.UserID,
		post.Content,
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

func (r *PostRepositoryImpl) List(ctx context.Context, limit int, offset int) ([]*domain.Post, error) {
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

	posts := make([]*domain.Post, 0)

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

		posts = append(posts, &post)
	}

	return posts, nil
}
