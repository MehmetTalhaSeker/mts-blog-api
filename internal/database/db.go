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

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(opts ...OptsFunc) (*PostgresStore, error) {
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

	return &PostgresStore{db: db}, nil
}

func (p PostgresStore) Init() error {
	return p.createUsersTable()
}

func (p PostgresStore) createUsersTable() error {
	query := `create table if not exists users (
    id serial primary key,
    encrypted_password varchar(500), 
    created_at timestamp,
    updated_at timestamp
	)`

	_, err := p.db.Exec(query)

	return err
}
