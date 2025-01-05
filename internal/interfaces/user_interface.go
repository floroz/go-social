package interfaces

import (
	"context"

	"github.com/floroz/go-social/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, createUser *domain.UserDTO) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, updateUser *domain.UserDTO) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]domain.User, error)
}

type UserService interface {
	Create(ctx context.Context, createUser *domain.UserDTO) (*domain.User, error)
	Update(ctx context.Context, updateUser *domain.UserDTO) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]domain.User, error)
}
