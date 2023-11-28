package comment

import (
	"context"
	"time"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/appcontext"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/rbac"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/repository"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type Service interface {
	Create(context.Context, *dto.CommentCreateRequest) error
	ReadsByPostID(p *pagination.Pageable, pid string) ([]*dto.CommentResponse, error)
	Delete(context.Context, *dto.RequestWithID) (*dto.ResponseWithID, error)
}

type service struct {
	repository repository.Comment
	rbac       rbac.RBAC
}

func NewService(rbac rbac.RBAC, repository repository.Comment) Service {
	return &service{
		repository: repository,
		rbac:       rbac,
	}
}

func (s *service) Create(ctx context.Context, req *dto.CommentCreateRequest) error {
	pid, err := apputils.StringToUINT64(req.PostID)
	if err != nil {
		return errorutils.New(errorutils.ErrInvalidID, err)
	}

	u, err := appcontext.MtsBlogUser(ctx)
	if err != nil {
		return err
	}

	var c model.Comment

	c.Author = u.Username
	c.CreatedAt = time.Now()
	c.Status = types.Active
	c.PostID = *pid
	c.UserID = u.UID
	c.Text = req.Text

	err = s.repository.Create(&c)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ReadsByPostID(p *pagination.Pageable, pid string) ([]*dto.CommentResponse, error) {
	comments, err := s.repository.ReadsByPostID(p, pid)
	if err != nil {
		return nil, err
	}

	var crs []*dto.CommentResponse

	for _, c := range *comments {
		crs = append(crs, c.ToDTO())
	}

	return crs, nil
}

func (s *service) Delete(ctx context.Context, req *dto.RequestWithID) (*dto.ResponseWithID, error) {
	cid, err := apputils.StringToUINT64(req.ID)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidID, err)
	}

	c, err := s.repository.Read(*cid)
	if err != nil {
		return nil, err
	}

	if !s.rbac.IsMe(ctx, c.UserID) && !s.rbac.IsModAuthorized(ctx) {
		return nil, errorutils.New(errorutils.ErrUnauthorized, nil)
	}

	if err = s.repository.Delete(*cid); err != nil {
		return nil, err
	}

	return &dto.ResponseWithID{ID: req.ID}, nil
}
