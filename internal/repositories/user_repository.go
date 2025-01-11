package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/interfaces"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) interfaces.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, createUser *domain.CreateUserDTO) (*domain.User, error) {

	user := domain.User{}

	query := `
        INSERT INTO users (first_name, last_name, email, username, password)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, first_name, last_name, email, username, password, created_at, updated_at
		`

	err := r.db.QueryRowContext(
		ctx,
		query,
		createUser.FirstName,
		createUser.LastName,
		createUser.Email,
		createUser.Username,
		createUser.Password,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return &user, err
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, userId int64) (*domain.User, error) {
	query := `
			SELECT id, first_name, last_name, email, username, password, created_at, updated_at
			FROM users
			WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, userId).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, username, password, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, username, password, created_at, updated_at
		FROM users
		WHERE username = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, userId int64, updateUser *domain.UpdateUserDTO) (*domain.User, error) {
	query := `
			UPDATE users
			SET first_name = $1, last_name = $2, email = $3, username = $4
			WHERE id = $5
			RETURNING id, first_name, last_name, email, username, password, created_at, updated_at
			`

	user := domain.User{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		updateUser.FirstName,
		updateUser.LastName,
		updateUser.Email,
		updateUser.Username,
		userId,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, userId int64) error {
	query := `
			DELETE FROM users
			WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, userId)

	return err
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	if limit == 0 {
		limit = 100
	}

	query := `
			SELECT id, first_name, last_name, email, username, password, created_at, updated_at
			FROM users
			LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		user := domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
