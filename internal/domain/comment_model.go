package domain

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EditableCommentFields struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

type CreateCommentDTO struct {
	EditableCommentFields
}

type UpdateCommentDTO struct {
	ID int `json:"id" validate:"required,min=1"`
	EditableCommentFields
}
