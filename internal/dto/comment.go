package dto

import (
	"time"
)

// CommentCreateRequest is the request body for the comment create endpoint.
type CommentCreateRequest struct {
	Text   string `json:"text"   validate:"required,min=2,max=100"`
	PostID string `json:"post_id" validate:"required"`
}

// CommentResponse is the response body for the comment.
type CommentResponse struct {
	Author    string     `json:"author"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	CreatedBy string     `json:"createdBy,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	ID        uint64     `json:"id,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
	UpdatedBy string     `json:"updatedBy,omitempty"`
	Text      string     `json:"text,omitempty"`
	PostID    uint64     `json:"post_id"`
	UserID    uint64     `json:"user_id"`
}

type ByPostIDRequest struct {
	PostID string `param:"pid" validate:"required"`
}
