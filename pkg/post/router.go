package post

import (
	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/rbac"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/repository"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
)

type Router struct {
	Authenticate   echo.MiddlewareFunc
	RBAC           rbac.RBAC
	RouterGroup    *echo.Group
	PostRepository repository.Post
}

func (r *Router) New() {
	ps := NewService(r.PostRepository)
	ph := NewHandler(ps)

	pgr := r.RouterGroup.Group("/posts")

	pgr.POST("", ph.Create(), r.Authenticate, r.RBAC.HasRole(types.Mod))
	pgr.GET("/:id", ph.Read())
	pgr.GET("", ph.Reads())
	pgr.PUT("/:id", ph.Update(), r.Authenticate, r.RBAC.HasRole(types.Mod))
	pgr.DELETE("/:id", ph.Delete(), r.Authenticate, r.RBAC.HasRole(types.Admin))
}
