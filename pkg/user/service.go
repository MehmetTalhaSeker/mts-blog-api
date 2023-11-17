package user

import (
	"context"
	"time"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/rbac"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type Service interface {
	Create(*dto.UserCreateRequest) error
	Read(*dto.RequestWithID) (*dto.UserResponse, error)
	Reads(*pagination.Pageable) ([]*dto.UserResponse, error)
	Update(context.Context, *dto.UserUpdateRequest) (*dto.UserResponse, error)
	Delete(*dto.RequestWithID) (*dto.ResponseWithID, error)
}

type service struct {
	repository Repository
	rbac       rbac.RBAC
}

func NewService(rbac rbac.RBAC, repository Repository) Service {
	return &service{
		repository: repository,
		rbac:       rbac,
	}
}

func (s *service) Create(req *dto.UserCreateRequest) error {
	var u model.User

	ep, err := apputils.EncryptPassword(req.Password)
	if err != nil {
		return errorutils.New(errorutils.ErrUnexpected, err)
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
		return err
	}

	return nil
}

func (s *service) Read(req *dto.RequestWithID) (*dto.UserResponse, error) {
	uid, err := apputils.StringToUINT64(req.ID)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidID, err)
	}

	u, err := s.repository.Read(*uid)
	if err != nil {
		return nil, err
	}

	return u.ToDTO(), nil
}

func (s *service) Reads(p *pagination.Pageable) ([]*dto.UserResponse, error) {
	users, err := s.repository.Reads(p)
	if err != nil {
		return nil, err
	}

	var usr []*dto.UserResponse

	for _, u := range *users {
		usr = append(usr, u.ToDTO())
	}

	return usr, nil
}

func (s *service) Update(ctx context.Context, req *dto.UserUpdateRequest) (*dto.UserResponse, error) {
	uid, err := apputils.StringToUINT64(req.ID)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidID, err)
	}

	if !s.rbac.IsMe(ctx, *uid) && !s.rbac.IsAdminAuthorized(ctx) {
		return nil, errorutils.New(errorutils.ErrUnauthorized, nil)
	}

	u, err := s.repository.Read(*uid)
	if err != nil {
		return nil, err
	}

	u.UpdatedAt = time.Now()
	u.Username = req.Username

	if err = s.repository.Update(u); err != nil {
		return nil, err
	}

	return u.ToDTO(), nil
}

func (s *service) Delete(req *dto.RequestWithID) (*dto.ResponseWithID, error) {
	uid, err := apputils.StringToUINT64(req.ID)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidID, err)
	}

	_, err = s.repository.Read(*uid)
	if err != nil {
		return nil, err
	}

	if err = s.repository.Delete(*uid); err != nil {
		return nil, err
	}

	return &dto.ResponseWithID{ID: req.ID}, nil
}
