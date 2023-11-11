package model

import (
	"time"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
)

type BaseModel struct {
	ID        uint64       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt *time.Time   `json:"deleted_at"`
	Status    types.Status `json:"status"`
	CreatedBy string       `json:"created_by"`
	UpdatedBy string       `json:"updated_by"`
	DeletedBy string       `json:"deleted_by"`
}
