package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID        int64      `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Username  string     `json:"username"`
	Password  string     `json:"-"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}

type EditableUserField struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=50"`
	LastName  string `json:"last_name" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,min=3,max=50,email"`
	Username  string `json:"username" validate:"required,min=3,max=50,alphanum"`
}

type CreateUserDTO struct {
	EditableUserField
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type UpdateUserDTO struct {
	EditableUserField
}

type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,min=3,max=50,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type UserClaims struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	jwt.StandardClaims
}
