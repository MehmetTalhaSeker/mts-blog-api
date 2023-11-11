package user

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

type Router struct {
	RouterGroup *echo.Group
	DB          *sql.DB
}

func (r *Router) New() {
	ur := NewRepository(r.DB)
	us := NewService(ur)
	uh := NewHandler(us)

	ugr := r.RouterGroup.Group("/users")

	ugr.POST("", uh.Create())
}
