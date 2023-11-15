package dto

import (
	"time"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
)

// UserCreateRequest is the request body for the user create endpoint.
type UserCreateRequest struct {
	Email    string `json:"email"          validate:"required,email"`
	Username string `json:"username"       validate:"required,min=3,max=21"`
	Password string `json:"password"       validate:"required,min=6,max=55"`
}

// UserUpdateRequest is the request body for the user update endpoint.
type UserUpdateRequest struct {
	ID       string `param:"id"         validate:"required"`
	Username string `json:"username"    validate:"required,min=3,max=21"`
}

// UserResponse is the response body for the user.
type UserResponse struct {
	CreatedAt      time.Time    `json:"createdAt,omitempty"`
	CreatedBy      string       `json:"createdBy,omitempty"`
	DeletedAt      *time.Time   `json:"deletedAt,omitempty"`
	Email          string       `json:"email,omitempty"`
	ID             uint64       `json:"id,omitempty"`
	Role           types.Role   `json:"role,omitempty"`
	Status         types.Status `json:"status,omitempty"`
	TermsOfService bool         `json:"termsOfService,omitempty"`
	UpdatedAt      time.Time    `json:"updatedAt,omitempty"`
	UpdatedBy      string       `json:"updatedBy,omitempty"`
	Username       string       `json:"username,omitempty"`
}
