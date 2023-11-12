package auth

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

type Router struct {
	RouterGroup *echo.Group
	DB          *sql.DB
}

func (r *Router) New() {
	ar := NewRepository(r.DB)
	as := NewService(ar)
	ah := NewHandler(as)

	ugr := r.RouterGroup.Group("/auth")

	ugr.POST("/login", ah.Login())
}
