package repository

import (
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
)

type Comment interface {
	Create(*model.Comment) error
	Read(id uint64) (*model.Comment, error)
	ReadsByPostID(p *pagination.Pageable, pid string) (*[]model.Comment, error)
	Delete(id uint64) error
}
