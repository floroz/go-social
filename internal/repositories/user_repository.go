package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"

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
        VALUES ($1, $2, $3, $4)
        RETURNING id, first_name, last_name, email, username`

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
	)

	if err != nil {
		log.Printf("error while creating user: %v", err)
		return nil, domain.NewInternalServerError("failed to create user")
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int) (*domain.User, error) {
	return nil, nil
	// 	query := `
	// 		SELECT id, first_name, last_name, email, password
	// 		FROM users
	// 		WHERE id = $1`

	// 	user := &domain.User{}
	// 	err := r.db.QueryRowContext(ctx, query, id).Scan(
	// 		&user.ID,
	// 		&user.FirstName,
	// 		&user.LastName,
	// 		&user.Email,
	// 		&user.Password,
	// 	)

	// 	if err != nil {
	// 		if errors.Is(err, sql.ErrNoRows) {
	// 			return nil, ErrNotFound
	// 		}
	// 		return nil, fmt.Errorf("get user by id: %w", err)
	// 	}

	// return user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, username
		FROM users
		WHERE email = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Printf("error while getting user by email: %v", err)
		return nil, domain.NewInternalServerError("failed to get user by email")
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, username
		FROM users
		WHERE username = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Printf("error while getting user by username: %v", err)
		return nil, domain.NewInternalServerError("failed to get user by username")
	}

	return user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, updateUser *domain.UpdateUserDTO) (*domain.User, error) {
	return nil, nil
	// 	query := `
	// 		UPDATE users
	// 		SET first_name = $1, last_name = $2, email = $3, password = $4
	// 		WHERE id = $5
	// 		RETURNING id, first_name, last_name, email, password
	// 		`

	// 	user := domain.User{}

	// 	err := r.db.QueryRowContext(
	// 		ctx,
	// 		query,
	// 		user.FirstName,
	// 		user.LastName,
	// 		user.Email,
	// 		user.Password,
	// 		user.ID,
	// 	).Scan(
	// 		&user.ID,
	// 		&user.FirstName,
	// 		&user.LastName,
	// 		&user.Email,
	// 		&user.Password,
	// 	)

	// 	if err != nil {
	// 		if errors.Is(err, sql.ErrNoRows) {
	// 			return nil, ErrNotFound
	// 		}
	// 		return nil, fmt.Errorf("update user: %w", err)
	// 	}

	// return &user, nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return nil, nil
	// 	query := `
	// 		SELECT id, first_name, last_name, email, password
	// 		FROM users
	// 		LIMIT $1 OFFSET $2`

	// 	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("list users: %w", err)
	// 	}
	// 	defer rows.Close()

	// 	users := []domain.User{}
	// 	for rows.Next() {
	// 		user := domain.User{}
	// 		err := rows.Scan(
	// 			&user.ID,
	// 			&user.FirstName,
	// 			&user.LastName,
	// 			&user.Email,
	// 			&user.Password,
	// 		)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("list users: %w", err)
	// 		}
	// 		users = append(users, user)
	// 	}

	// return users, nil
}
