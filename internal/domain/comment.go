package domain

import "time"

type Comment struct {
	ID              int64     `json:"id"`
	PostID          int64     `json:"post_id"`
	UserID          int64     `json:"user_id"`
	ParentCommentID *int64    `json:"parent_comment_id,omitempty"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CommentDTO struct {
	PostID          int64  `json:"post_id"`
	UserID          int64  `json:"user_id"`
	ParentCommentID *int64 `json:"parent_comment_id,omitempty"`
	Content         string `json:"content"`
}
