package repository

import (
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
)

type Post interface {
	Create(*model.Post) error
	Read(id uint64) (*model.Post, error)
	Reads(*pagination.Pageable) (*[]model.Post, error)
	Update(*model.Post) error
	Delete(id uint64) error
}
