package interfaces

import (
	"context"
	"time"

	"github.com/floroz/go-social/internal/domain"
)

type AuthService interface {
	GenerateJWTToken(user *domain.User, expiration time.Duration) (string, error)
	Login(ctx context.Context, loginUser *domain.LoginUserDTO) (*domain.User, error)
}
