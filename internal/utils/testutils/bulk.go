package testutils

import (
	"database/sql"
	"log"

	"github.com/lib/pq"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
)

func InsertUsers(us []*model.User, db *sql.DB) {
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("users", "id", "email", "username", "encrypted_password", "user_role", "created_at", "updated_at"))
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range us {
		_, err = stmt.Exec(user.ID, user.Email, user.Username, user.EncryptedPassword, user.Role, user.CreatedAt, user.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
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
