package domain

import "time"

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	ID int `json:"id" validate:"required,min=1"`
	EditableUserField
}
