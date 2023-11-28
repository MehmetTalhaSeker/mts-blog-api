package post

import (
	"time"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/repository"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type Service interface {
	Create(*dto.PostCreateRequest) error
	Read(*dto.RequestWithID) (*dto.PostResponse, error)
	Reads(*pagination.Pageable) ([]*dto.PostResponse, error)
	Update(*dto.PostUpdateRequest) (*dto.PostResponse, error)
	Delete(*dto.RequestWithID) (*dto.ResponseWithID, error)
}

type service struct {
	repository repository.Post
}

func NewService(repository repository.Post) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(req *dto.PostCreateRequest) error {
	var u model.Post

	u.Body = req.Body
	u.CreatedAt = time.Now()
	u.Status = types.Active
	u.UpdatedAt = time.Now()
	u.Title = req.Title

	err := s.repository.Create(&u)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Read(req *dto.RequestWithID) (*dto.PostResponse, error) {
	pid, err := apputils.StringToUINT64(req.ID)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidID, err)
	}

	p, err := s.repository.Read(*pid)
	if err != nil {
		return nil, err
	}

	return p.ToDTO(), nil
}

func (s *service) Reads(p *pagination.Pageable) ([]*dto.PostResponse, error) {
	posts, err := s.repository.Reads(p)
	if err != nil {
		return nil, err
	}

	var psr []*dto.PostResponse

	for _, u := range *posts {
		psr = append(psr, u.ToDTO())
	}

	return psr, nil
}

func (s *service) Update(req *dto.PostUpdateRequest) (*dto.PostResponse, error) {
	uid, err := apputils.StringToUINT64(req.ID)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidID, err)
	}

	p, err := s.repository.Read(*uid)
	if err != nil {
		return nil, err
	}

	if p.Title == req.Title && req.Body == "" {
		return nil, nil
	}

	if req.Body != "" {
		p.Body = req.Body
	}

	p.Title = req.Title
	p.UpdatedAt = time.Now()

	if err = s.repository.Update(p); err != nil {
		return nil, err
	}

	return p.ToDTO(), nil
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
