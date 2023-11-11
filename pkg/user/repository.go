package user

import (
	"database/sql"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
)

type Repository interface {
	Create(u *model.User) error
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
	query := `insert into users 
    (email, username, encrypted_password,created_at, updated_at)
    values ($1, $2, $3, $4, $5)`

	_, err := r.db.Query(query, u.Email, u.Username, u.EncryptedPassword, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}

	return err
}
