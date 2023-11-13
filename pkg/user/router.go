package user

import (
	"database/sql"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/rbac"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"

	"github.com/labstack/echo/v4"
)

type Router struct {
	Authenticate echo.MiddlewareFunc
	DB           *sql.DB
	RBAC         rbac.RBAC
	RouterGroup  *echo.Group
}

func (r *Router) New() {
	ur := NewRepository(r.DB)
	us := NewService(ur)
	uh := NewHandler(us)

	ugr := r.RouterGroup.Group("/users", r.Authenticate)

	ugr.POST("", uh.Create(), r.RBAC.HasRole(types.Admin))
}
