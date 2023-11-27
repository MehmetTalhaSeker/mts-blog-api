package auth

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/repository"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type Service interface {
	Login(*dto.LoginRequest) (*dto.WithTokenResponse, error)
	Register(*dto.RegisterRequest) (*dto.WithTokenResponse, error)
}

type service struct {
	userRepository repository.User
}

func NewService(repository repository.User) Service {
	return &service{
		userRepository: repository,
	}
}

func (s *service) Login(req *dto.LoginRequest) (*dto.WithTokenResponse, error) {
	u, err := s.userRepository.ReadByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(req.Password)); err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidPassword, err)
	}

	token, err := apputils.CreateJWT(u)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrUnexpected, err)
	}

	return &dto.WithTokenResponse{
		Token: token,
		Claims: dto.Claims{
			UID:      u.ID,
			Role:     u.Role,
			Username: u.Username,
			Email:    u.Email,
		},
	}, nil
}

func (s *service) Register(req *dto.RegisterRequest) (*dto.WithTokenResponse, error) {
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

	err = s.userRepository.Create(&u)
	if err != nil {
		return nil, err
	}

	token, err := apputils.CreateJWT(&u)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrUnexpected, err)
	}

	return &dto.WithTokenResponse{
		Token: token,
		Claims: dto.Claims{
			Role:     u.Role,
			Username: u.Username,
			Email:    u.Email,
		},
	}, nil
}
