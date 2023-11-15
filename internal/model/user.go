package model

import (
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
)

type User struct {
	BaseModel
	Email             string     `json:"email"`
	Role              types.Role `json:"role"`
	Username          string     `json:"username"`
	EncryptedPassword string     `json:"-"`
}

func (u User) ToDTO() *dto.UserResponse {
	return &dto.UserResponse{
		CreatedAt: u.CreatedAt,
		CreatedBy: u.CreatedBy,
		DeletedAt: u.DeletedAt,
		Email:     u.Email,
		ID:        u.ID,
		Role:      u.Role,
		Status:    u.Status,
		UpdatedAt: u.UpdatedAt,
		UpdatedBy: u.UpdatedBy,
		Username:  u.Username,
	}
}
