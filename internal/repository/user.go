package repository

import (
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
)

type User interface {
	Create(*model.User) error
	Read(id uint64) (*model.User, error)
	ReadByEmail(email string) (*model.User, error)
	Reads(*pagination.Pageable) (*[]model.User, error)
	Update(*model.User) error
	Delete(id uint64) error
}
