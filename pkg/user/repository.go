package user

import (
	"database/sql"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/dbutils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type Repository interface {
	Create(user *model.User) error
	Read(id uint64) (*model.User, error)
	ReadByEmail(email string) (*model.User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(u *model.User) error {
	query := `INSERT INTO users 
    (email, username, encrypted_password, user_role, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Query(query, u.Email, u.Username, u.EncryptedPassword, u.Role, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return errorutils.New(errorutils.ErrUserCreate, err)
	}

	return nil
}

func (r *repository) Read(i uint64) (*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE id = $1", i)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidRequest, err)
	}

	for rows.Next() {
		user, err := dbutils.ScanIntoUser(rows)
		if err != nil {
			return nil, errorutils.New(errorutils.ErrUserNotFound, err)
		}

		return user, nil
	}

	return nil, errorutils.New(errorutils.ErrUserNotFound, nil)
}

func (r *repository) ReadByEmail(e string) (*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE email = $1", e)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidRequest, err)
	}

	for rows.Next() {
		user, err := dbutils.ScanIntoUser(rows)
		if err != nil {
			return nil, errorutils.New(errorutils.ErrUserNotFound, err)
		}

		return user, nil
	}

	return nil, errorutils.New(errorutils.ErrEmailNotFound, errorutils.ErrUserRead)
}
