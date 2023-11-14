package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Opts struct {
	user     string
	name     string
	password string
	port     string
}

type OptsFunc func(*Opts)

func defaultOpts() Opts {
	return Opts{
		user:     "development",
		name:     "development",
		password: "development",
		port:     "5432",
	}
}

func WithUser(u string) OptsFunc {
	return func(opts *Opts) {
		opts.user = u
	}
}

func WithPassword(p string) OptsFunc {
	return func(opts *Opts) {
		opts.password = p
	}
}

func WithName(n string) OptsFunc {
	return func(opts *Opts) {
		opts.name = n
	}
}

func WithPort(p string) OptsFunc {
	return func(opts *Opts) {
		opts.port = p
	}
}

type Store struct {
	DB *sql.DB
}

func NewPostgresStore(opts ...OptsFunc) (*Store, error) {
	o := defaultOpts()
	for _, fn := range opts {
		fn(&o)
	}

	conStr := fmt.Sprintf(`user=%s dbname=%s password=%s port=%s sslmode=disable`, o.user, o.name, o.password, o.port)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{DB: db}, nil
}

func (s Store) Init() error {
	if err := s.createUsersEnumRoles(); err != nil {
		return err
	}

	if err := s.createUsersTable(); err != nil {
		return err
	}

	return nil
}

func (s Store) createUsersEnumRoles() error {
	query := `DO $$ BEGIN
	IF to_regtype('user_roles') IS NULL THEN
	CREATE TYPE user_roles AS ENUM('admin', 'mod', 'registered');
	END IF;
	END $$;`

	_, err := s.DB.Exec(query)

	return err
}

func (s Store) createUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
    id serial 		   primary key,
    encrypted_password varchar(500) NOT NULL, 
    username 		   varchar(21) NOT NULL UNIQUE, 
    email 			   varchar(55) NOT NULL UNIQUE,
	user_role     	   user_roles,
    created_at 		   timestamp,
    updated_at 		   timestamp
	)`

	_, err := s.DB.Exec(query)

	return err
}
