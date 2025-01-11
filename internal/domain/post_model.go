package domain

import (
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

type EditablePostFields struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

type CreatePostDTO struct {
	EditablePostFields
}

type UpdatePostDTO struct {
	EditablePostFields
}
