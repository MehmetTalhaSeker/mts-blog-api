package dto

import (
	"time"
)

// PostCreateRequest is the request body for the post create endpoint.
type PostCreateRequest struct {
	Title string `json:"title"          validate:"required,min=3,max=21"`
	Body  string `json:"body"       validate:"required,min=1"`
}

// PostUpdateRequest is the request body for the post update endpoint.
type PostUpdateRequest struct {
	ID    string `param:"id"         validate:"required"`
	Title string `json:"title"    validate:"required,min=3,max=21"`
	Body  string `json:"body"       validate:"omitempty,min=1"`
}

// PostResponse is the response body for the post.
type PostResponse struct {
	Body      string     `json:"body,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	CreatedBy string     `json:"createdBy,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	ID        uint64     `json:"id,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
	UpdatedBy string     `json:"updatedBy,omitempty"`
	Title     string     `json:"title,omitempty"`
}
