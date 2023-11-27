package main

import (
	"database/sql"
	"log"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/database"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/rbac"
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

	// Create a Postgres store.
	store := database.NewPostgresStore(database.WithUser(cfg.DB.User), database.WithName(cfg.DB.Name), database.WithPassword(cfg.DB.Password))

	// Initialize the Postgres store.
	store.InitDB()

	// Initialize Role based access control.
	rb := rbac.New()

	app := &application{
		config: cfg,
		db:     store.GetInstance(),
		rbac:   rb,
	}

	log.Printf("starting server on %s:%s (version %s)", cfg.Rest.Host, cfg.Rest.Port, cfg.Rest.Version)
	app.start()
}
