package domain

import "time"

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EditableCommentFields struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

type CreateCommentDTO struct {
	PostID int `json:"post_id" validate:"required"`
	UserID int `json:"user_id" validate:"required"`
	EditableCommentFields
}

type UpdateCommentDTO struct {
	ID int `json:"id" validate:"required"`
	EditableCommentFields
}
