package main

import (
	"database/sql"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/rbac"
	"log"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/database"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/config"
)

type application struct {
	config *config.Config
	db     *sql.DB
	rbac   rbac.RBAC
}

func main() {
	// Initialize application configs.
	cfg := config.Init()

	// Create a Postgres Store.
	store, err := database.NewPostgresStore(database.WithUser(cfg.DB.User), database.WithName(cfg.DB.Name), database.WithPassword(cfg.DB.Password))
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the Postgres Store.
	err = store.Init()
	if err != nil {
		log.Fatal(err)
	}

	rb := rbac.New()

	app := &application{
		config: cfg,
		db:     store.DB,
		rbac:   rb,
	}

	log.Printf("starting server on %s:%s (version %s)", cfg.Rest.Host, cfg.Rest.Port, cfg.Rest.Version)
	app.start()
}
