package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	singleton *postgresStore
	storeOnce sync.Once
	initOnce  sync.Once
)

type postgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(opts ...StoreOptsFunc) SQLStore {
	storeOnce.Do(func() {
		o := postgresStoreDefaultOpts()
		for _, fn := range opts {
			fn(&o)
		}

		conStr := fmt.Sprintf(`user=%s dbname=%s password=%s port=%s sslmode=disable`, o.user, o.name, o.password, o.port)

		db, err := sql.Open("postgres", conStr)
		if err != nil {
			log.Fatal(err.Error())
		}

		singleton = &postgresStore{DB: db}

		if err := db.Ping(); err != nil {
			log.Fatal("Connection Error: " + err.Error())
		}
	})

	return singleton
}

func (s postgresStore) GetInstance() *sql.DB {
	return s.DB
}

func (s postgresStore) InitDB() {
	initOnce.Do(func() {
		s.createUsersEnumRoles()
		s.createUsersTable()
		s.createPostsTable()
		s.createCommentsTable()
	})
}

func (s postgresStore) createUsersEnumRoles() {
	query := `DO $$ BEGIN
	IF to_regtype('user_roles') IS NULL THEN
	CREATE TYPE user_roles AS ENUM('admin', 'mod', 'registered');
	END IF;
	END $$;
	`

	_, err := s.DB.Exec(query)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (s postgresStore) createUsersTable() {
	query := `CREATE TABLE IF NOT EXISTS users (
    id				   serial PRIMARY KEY,
    encrypted_password varchar(500) NOT NULL, 
    username 		   varchar(21) NOT NULL UNIQUE, 
    email 			   varchar(55) NOT NULL UNIQUE,
	user_role     	   user_roles,
    created_at 		   timestamp,
    updated_at 		   timestamp
	)`

	_, err := s.DB.Exec(query)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (s postgresStore) createPostsTable() {
	query := `CREATE TABLE IF NOT EXISTS posts (
    id 				   serial PRIMARY KEY,
    title 			   varchar(255),
	body 			   varchar, 
    created_at 		   timestamp,
    updated_at 		   timestamp
	)`

	_, err := s.DB.Exec(query)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (s postgresStore) createCommentsTable() {
	query := `CREATE TABLE IF NOT EXISTS comments (
    id 				   serial PRIMARY KEY,
    text 			   varchar(255),
	author 			   varchar references users(username), 
	user_id 		   int references users(id),
	post_id 		   int references posts(id),
    created_at 		   timestamp,
    updated_at 		   timestamp
	)`

	_, err := s.DB.Exec(query)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func postgresStoreDefaultOpts() StoreOpts {
	return StoreOpts{
		user:     "development",
		name:     "development",
		password: "development",
		port:     "5432",
	}
}
