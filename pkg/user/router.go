package user

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/rbac"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
)

type Router struct {
	Authenticate echo.MiddlewareFunc
	DB           *sql.DB
	RBAC         rbac.RBAC
	RouterGroup  *echo.Group
}

func (r *Router) New() {
	ur := NewRepository(r.DB)
	us := NewService(r.RBAC, ur)
	uh := NewHandler(us)

	ugr := r.RouterGroup.Group("/users", r.Authenticate)

	ugr.POST("", uh.Create(), r.RBAC.HasRole(types.Admin))
	ugr.GET("/:id", uh.Read(), r.RBAC.HasRole(types.Mod))
	ugr.PUT("/:id", uh.Update(), r.RBAC.HasRole(types.Registered))
	ugr.DELETE("/:id", uh.Delete(), r.RBAC.HasRole(types.Admin))

	ugr.GET("", uh.Read(), r.RBAC.HasRole(types.Registered))
}
