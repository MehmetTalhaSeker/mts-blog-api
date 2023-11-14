package dbutils

import (
	"database/sql"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
)

func ScanIntoUser(rows *sql.Rows) (*model.User, error) {
	u := new(model.User)
	err := rows.Scan(&u.ID, &u.EncryptedPassword, &u.Username, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)

	return u, err
}
