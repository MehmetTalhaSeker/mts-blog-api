package auth

import (
	"database/sql"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type Repository interface {
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

func (r repository) Read(id uint64) (*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidRequest, err)
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, errorutils.New(errorutils.ErrUserNotFound, nil)
}

func (r repository) ReadByEmail(email string) (*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidRequest, err)
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, errorutils.New(errorutils.ErrEmailNotFound, errorutils.ErrUserRead)
}

func scanIntoUser(rows *sql.Rows) (*model.User, error) {
	u := new(model.User)
	err := rows.Scan(&u.ID, &u.EncryptedPassword, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt)

	return u, err
}
