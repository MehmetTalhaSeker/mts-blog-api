package dto

import "github.com/MehmetTalhaSeker/mts-blog-api/internal/types"

// Claims is the JWT claims.
type Claims struct {
	UID      uint64     `json:"uid"`
	Role     types.Role `json:"role"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
}

// LoginRequest is the request body for the user login endpoint.
type LoginRequest struct {
	Email    string `json:"email"          validate:"required,email"`
	Password string `json:"password"       validate:"required,min=6,max=55"`
}

// LoginResponse is the response body for the user login endpoint.
type LoginResponse struct {
	Token string `json:"token"`
}
