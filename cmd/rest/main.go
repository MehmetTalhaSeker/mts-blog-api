package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/database"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/config"
)

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

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
