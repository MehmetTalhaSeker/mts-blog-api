package auth

import (
	"os"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type Service interface {
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	u, err := s.repository.ReadByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(req.Password)); err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidPassword, err)
	}

	token, err := createJWT(u)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrUnexpected, err)
	}

	resp := dto.LoginResponse{
		Token: token,
	}

	return &resp, nil
}

func createJWT(u *model.User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"email":     u.Email,
		"username":  u.Username,
		"role":      u.Role,
		"uid":       u.ID,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
