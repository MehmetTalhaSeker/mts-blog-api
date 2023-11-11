package model

import "github.com/MehmetTalhaSeker/mts-blog-api/internal/types"

type User struct {
	BaseModel
	Email             string     `json:"email"`
	Role              types.Role `json:"role"`
	Username          string     `json:"username"`
	EncryptedPassword string     `json:"-"`
}
