package domain

import "time"

type Comment struct {
	ID              int       `json:"id"`
	PostID          int       `json:"post_id"`
	UserID          int       `json:"user_id"`
	ParentCommentID *int      `json:"parent_comment_id,omitempty"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateCommentDTO struct {
	PostID          int    `json:"post_id" validate:"required"`
	UserID          int    `json:"user_id" validate:"required"`
	ParentCommentID *int   `json:"parent_comment_id,omitempty"`
	Content         string `json:"content" validate:"required,min=1,max=1000"`
}
