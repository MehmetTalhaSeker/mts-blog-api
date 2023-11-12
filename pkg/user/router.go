package user

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

type Router struct {
	Authenticate echo.MiddlewareFunc
	DB           *sql.DB
	RouterGroup  *echo.Group
}

func (r *Router) New() {
	ur := NewRepository(r.DB)
	us := NewService(ur)
	uh := NewHandler(us)

	ugr := r.RouterGroup.Group("/users", r.Authenticate)

	ugr.POST("", uh.Create())
}
