package user

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
)

type Service interface {
	Create(ctx context.Context, req *dto.UserCreateRequest) (*uint64, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(ctx context.Context, req *dto.UserCreateRequest) (*uint64, error) {
	var u model.User

	cp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.CreatedAt = time.Now()
	u.EncryptedPassword = string(cp)
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
