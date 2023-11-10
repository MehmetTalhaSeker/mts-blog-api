package model

import (
	"time"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
)

type Model struct {
	ID        uint64       `json:"id"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	DeletedAt *time.Time   `json:"deletedAt"`
	Status    types.Status `json:"status"`
	CreatedBy string       `json:"createdBy"`
	UpdatedBy string       `json:"updatedBy"`
	DeletedBy string       `json:"deletedBy"`
}
