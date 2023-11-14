package auth

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/pkg/user"
)

type Router struct {
	RouterGroup *echo.Group
	DB          *sql.DB
}

func (r *Router) New() {
	ur := user.NewRepository(r.DB)
	as := NewService(ur)
	ah := NewHandler(as)

	ugr := r.RouterGroup.Group("/auth")

	ugr.POST("/login", ah.Login())
	ugr.POST("/register", ah.Register())
}
