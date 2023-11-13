package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/validatorutils"
	"github.com/MehmetTalhaSeker/mts-blog-api/pkg/auth"
	"github.com/MehmetTalhaSeker/mts-blog-api/pkg/user"
)

func (app *application) start() {
	e := echo.New()
	e.Use(
		// middleware.Recover(), // Recover from all panics to always have your server up
		middleware.Logger(),    // Log everything to stdout
		middleware.RequestID(), // Generate a request id on the HTTP response headers for identification
		middleware.CORSWithConfig(middleware.DefaultCORSConfig),
	)

	e.Validator = validatorutils.NewValidator()
	e.HTTPErrorHandler = errorutils.Handler

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// create a new router group.
	routerGroup := e.Group("v1")

	// user router initialization.
	userRouter := &user.Router{
		Authenticate: app.authenticate(),
		DB:           app.db,
		RBAC:         app.rbac,
		RouterGroup:  routerGroup,
	}
	userRouter.New()

	// auth router initialization.
	authRouter := &auth.Router{
		DB:          app.db,
		RouterGroup: routerGroup,
	}
	authRouter.New()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", app.config.Rest.Port)))
}
