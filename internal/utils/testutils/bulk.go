package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
)

func InsertUsers(us []*model.User, db *sql.DB) {
	vs := make([]string, 0, len(us))
	va := make([]interface{}, 0, len(us)*3)

	for _, u := range us {
		vs = append(vs, "($1, $2, $3, $4, $5, $6)")
		va = append(va, u.Email)
		va = append(va, u.Username)
		va = append(va, u.EncryptedPassword)
		va = append(va, u.Role)
		va = append(va, u.CreatedAt)
		va = append(va, u.UpdatedAt)
	}

	stmt := fmt.Sprintf("INSERT INTO users (email, username, encrypted_password, user_role, created_at, updated_at) VALUES %s",
		strings.Join(vs, ","))

	_, err := db.Exec(stmt, va...)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func DeleteUsers(db *sql.DB) {
	dcq := "DROP TABLE comments;"

	_, err := db.Exec(dcq)
	if err != nil {
		log.Fatalf(err.Error())
	}

	tuq := "TRUNCATE TABLE users"

	_, err = db.Exec(tuq)
	if err != nil {
		log.Fatalf(err.Error())
	}

	ctc := `CREATE TABLE IF NOT EXISTS comments (
    id 				   serial PRIMARY KEY,
    text 			   varchar(255),
	author 			   varchar references users(username), 
	user_id 		   int references users(id),
	post_id 		   int references posts(id),
    created_at 		   timestamp,
    updated_at 		   timestamp
	)`

	_, err = db.Exec(ctc)
	if err != nil {
		log.Fatal(err.Error())
	}
}
