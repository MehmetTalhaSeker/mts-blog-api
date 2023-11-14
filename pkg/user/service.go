package user

import (
	"time"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type Service interface {
	Create(req *dto.UserCreateRequest) (*uint64, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(req *dto.UserCreateRequest) (*uint64, error) {
	var u model.User

	ep, err := apputils.EncryptPassword(req.Password)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrUnexpected, err)
	}

	u.CreatedAt = time.Now()
	u.EncryptedPassword = ep
	u.Email = req.Email
	u.Role = types.Registered
	u.Status = types.Active
	u.UpdatedAt = time.Now()
	u.Username = req.Username

	err = s.repository.Create(&u)
	if err != nil {
		return nil, err
	}

	return &u.ID, nil
}
