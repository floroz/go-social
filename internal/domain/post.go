package domain

import (
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EditablePostFields struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

type CreatePostDTO struct {
	UserID int `json:"user_id" validate:"required,min=1"`
	EditablePostFields
}

type UpdatePostDTO struct {
	ID int `json:"id" validate:"required,min=1"`
	EditablePostFields
}
