package services

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/env"
	"github.com/floroz/go-social/internal/interfaces"
	"github.com/floroz/go-social/internal/validation"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo interfaces.UserRepository
}

func NewAuthService(userRepo interfaces.UserRepository) *authService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) GenerateJWTToken(user *domain.User, expiration time.Duration) (string, error) {
	jwtSecret := env.GetJWTSecret()

	if jwtSecret == "" {
		log.Error().Msg("jwt secret is not set")
		return "", domain.NewInternalServerError("jwt secret is not set")
	}

	claims := domain.UserClaims{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		log.Error().Err(err).Msg("failed to generate jwt token")
		return "", domain.NewInternalServerError("failed to generate jwt token")
	}

	return signedToken, nil
}

func (s *authService) Login(ctx context.Context, loginUser *domain.LoginUserDTO) (*domain.User, error) {
	err := validation.Validate.Struct(loginUser)

	if err != nil {
		return nil, domain.NewBadRequestError(err.Error())
	}

	user, err := s.userRepo.GetByEmail(ctx, loginUser.Email)

	if err != nil && err == domain.ErrNotFound {
		log.Error().Err(err).Msg("attempting to login a user that doesn't exist")
		return nil, domain.NewUnauthorizedError("invalid email or password")
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to get user by email")
		return nil, domain.NewInternalServerError("failed to get user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		log.Error().Err(err).Msg("failed to compare password")
		return nil, domain.NewUnauthorizedError("invalid email or password")
	}

	return user, nil
}
