package postgresadapter

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/repository"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.User {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(u *model.User) error {
	query := `INSERT INTO users 
    (email, username, encrypted_password, user_role, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Query(query, u.Email, u.Username, u.EncryptedPassword, u.Role, u.CreatedAt, u.UpdatedAt)

	var pErr *pq.Error
	if err != nil {
		switch errors.As(err, &pErr) {
		case pErr.Constraint == "users_username_key":
			return errorutils.New(errorutils.ErrUsernameAlreadyTaken, err)
		case pErr.Constraint == "users_email_key":
			return errorutils.New(errorutils.ErrEmailAlreadyTaken, err)
		}

		return errorutils.New(errorutils.ErrUserCreate, err)
	}

	return nil
}

func (r *userRepository) Read(i uint64) (*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE id = $1", i)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidRequest, err)
	}

	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, errorutils.New(errorutils.ErrUserNotFound, err)
		}

		return user, nil
	}

	return nil, errorutils.New(errorutils.ErrUserNotFound, errorutils.ErrUserRead)
}

func (r *userRepository) ReadByEmail(e string) (*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE email = $1", e)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidRequest, err)
	}

	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, errorutils.New(errorutils.ErrUserNotFound, err)
		}

		return user, nil
	}

	return nil, errorutils.New(errorutils.ErrEmailNotFound, errorutils.ErrUserRead)
}

func (r *userRepository) Reads(p *pagination.Pageable) (*[]model.User, error) {
	// Note: Just for show off. I know it can be handled in single query :)
	fq := `SELECT * FROM users ORDER BY ` +
		fmt.Sprintf("%s LIMIT $1 OFFSET $2;", p.Order())

	cq := `SELECT COUNT(*) FROM users;`

	countErr := make(chan error)
	findErr := make(chan error)

	var users []model.User

	go func() {
		rows, err := r.db.Query(fq, p.Size, p.Offset())
		if err != nil {
			findErr <- errorutils.New(errorutils.ErrUserReads, err)
		}

		for rows.Next() {
			u, err := scanIntoUser(rows)
			if err != nil {
				findErr <- errorutils.New(errorutils.ErrUserReads, err)
			}

			users = append(users, *u)
		}
		findErr <- nil
	}()

	var count int64
	go func() {
		rows, err := r.db.Query(cq)
		if err != nil {
			countErr <- errorutils.New(errorutils.ErrUserCount, err)
		}

		for rows.Next() {
			err = rows.Scan(&count)
			if err != nil {
				countErr <- errorutils.New(errorutils.ErrUserCount, err)
			}
		}

		p.TotalCount = count
		countErr <- nil
	}()

	err := <-countErr
	if err != nil {
		return nil, err
	}

	err = <-findErr
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *userRepository) Update(u *model.User) error {
	_, err := r.db.Query("UPDATE users SET username = $1, updated_at = $2 WHERE id = $3;", u.Username, u.UpdatedAt, u.ID)

	var pErr *pq.Error

	if err != nil {
		switch errors.As(err, &pErr) {
		case pErr.Constraint == "users_username_key":
			return errorutils.New(errorutils.ErrUsernameAlreadyTaken, err)
		}

		return errorutils.New(errorutils.ErrUserUpdate, err)
	}

	return nil
}

func (r *userRepository) Delete(i uint64) error {
	_, err := r.db.Query("DELETE FROM users WHERE id = $1", i)
	if err != nil {
		return errorutils.New(errorutils.ErrUserDelete, err)
	}

	return nil
}

func scanIntoUser(rows *sql.Rows) (*model.User, error) {
	u := new(model.User)
	err := rows.Scan(&u.ID, &u.EncryptedPassword, &u.Username, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)

	return u, err
}
