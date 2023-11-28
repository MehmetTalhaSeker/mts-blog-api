package postgresadapter

import (
	"database/sql"
	"fmt"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/repository"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) repository.Comment {
	return &commentRepository{
		db: db,
	}
}

func (r *commentRepository) Create(c *model.Comment) error {
	query := `INSERT INTO comments 
    (author, post_id, text, user_id, created_at)
    VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Query(query, c.Author, c.PostID, c.Text, c.UserID, c.CreatedAt)
	if err != nil {
		return errorutils.New(errorutils.ErrCommentCreate, err)
	}

	return nil
}

func (r *commentRepository) Read(id uint64) (*model.Comment, error) {
	rows, err := r.db.Query("SELECT * FROM comments WHERE id = $1", id)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidRequest, err)
	}

	for rows.Next() {
		comment, err := scanIntoComment(rows)
		if err != nil {
			return nil, errorutils.New(errorutils.ErrCommentNotFound, err)
		}

		return comment, nil
	}

	return nil, errorutils.New(errorutils.ErrCommentNotFound, errorutils.ErrCommentRead)
}

func (r *commentRepository) ReadsByPostID(p *pagination.Pageable, pid string) (*[]model.Comment, error) {
	fq := `SELECT * FROM comments WHERE comments.post_id=$1 ORDER BY ` +
		fmt.Sprintf("%s LIMIT $2 OFFSET $3;", p.Order())

	cq := `SELECT COUNT(*) FROM comments WHERE post_id=$1;`

	countErr := make(chan error)
	findErr := make(chan error)

	var comments []model.Comment

	go func() {
		rows, err := r.db.Query(fq, pid, p.Size, p.Offset())
		if err != nil {
			findErr <- errorutils.New(errorutils.ErrCommentReads, err)
		}

		for rows.Next() {
			u, err := scanIntoComment(rows)
			if err != nil {
				findErr <- errorutils.New(errorutils.ErrCommentReads, err)
			}

			comments = append(comments, *u)
		}
		findErr <- nil
	}()

	var count int64
	go func() {
		rows, err := r.db.Query(cq, pid)
		if err != nil {
			countErr <- errorutils.New(errorutils.ErrCommentCount, err)
		}

		for rows.Next() {
			err = rows.Scan(&count)
			if err != nil {
				countErr <- errorutils.New(errorutils.ErrCommentCount, err)
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

	return &comments, nil
}

func (r *commentRepository) Delete(id uint64) error {
	_, err := r.db.Query("DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		return errorutils.New(errorutils.ErrCommentDelete, err)
	}

	return nil
}

func scanIntoComment(rows *sql.Rows) (*model.Comment, error) {
	c := new(model.Comment)
	err := rows.Scan(&c.ID, &c.Author, &c.UserID, &c.PostID, &c.Text, &c.CreatedAt)

	return c, err
}
