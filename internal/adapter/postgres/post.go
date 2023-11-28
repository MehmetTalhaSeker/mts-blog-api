package postgresadapter

import (
	"database/sql"
	"fmt"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/repository"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) repository.Post {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Create(p *model.Post) error {
	query := `INSERT INTO posts 
    (title, body, created_at, updated_at)
    VALUES ($1, $2, $3, $4)`

	_, err := r.db.Query(query, p.Title, p.Body, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return errorutils.New(errorutils.ErrPostCreate, err)
	}

	return nil
}

func (r *postRepository) Read(id uint64) (*model.Post, error) {
	rows, err := r.db.Query("SELECT * FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidRequest, err)
	}

	for rows.Next() {
		post, err := scanIntoPost(rows)
		if err != nil {
			return nil, errorutils.New(errorutils.ErrPostNotFound, err)
		}

		return post, nil
	}

	return nil, errorutils.New(errorutils.ErrPostNotFound, errorutils.ErrPostRead)
}

func (r *postRepository) Reads(p *pagination.Pageable) (*[]model.Post, error) {
	q := `SELECT id, title, body, created_at, updated_at, COUNT(*) OVER() AS count FROM posts ORDER BY ` +
		fmt.Sprintf("%s LIMIT $1 OFFSET $2;", p.Order())

	var posts []model.Post

	var count int64

	rows, err := r.db.Query(q, p.Size, p.Offset())
	if err != nil {
		return nil, errorutils.New(errorutils.ErrPostReads, err)
	}

	for rows.Next() {
		p := new(model.Post)

		err := rows.Scan(&p.ID, &p.Title, &p.Body, &p.CreatedAt, &p.UpdatedAt, &count)
		if err != nil {
			return nil, errorutils.New(errorutils.ErrPostReads, err)
		}

		posts = append(posts, *p)
	}

	p.TotalCount = count

	return &posts, nil
}

func (r *postRepository) Update(p *model.Post) error {
	_, err := r.db.Query("UPDATE posts SET title = $1, body = $2, updated_at = $3 WHERE id = $4;", p.Title, p.Body, p.UpdatedAt, p.ID)
	if err != nil {
		return errorutils.New(errorutils.ErrPostUpdate, err)
	}

	return nil
}

func (r *postRepository) Delete(id uint64) error {
	_, err := r.db.Query("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return errorutils.New(errorutils.ErrPostDelete, err)
	}

	return nil
}

func scanIntoPost(rows *sql.Rows) (*model.Post, error) {
	p := new(model.Post)
	err := rows.Scan(&p.ID, &p.Title, &p.Body, &p.CreatedAt, &p.UpdatedAt)

	return p, err
}
