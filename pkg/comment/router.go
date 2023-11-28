package comment

import (
	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/rbac"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/repository"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
)

type Router struct {
	Authenticate      echo.MiddlewareFunc
	RBAC              rbac.RBAC
	RouterGroup       *echo.Group
	CommentRepository repository.Comment
}

func (r *Router) New() {
	cs := NewService(r.RBAC, r.CommentRepository)
	ch := NewHandler(cs)

	cgr := r.RouterGroup.Group("/comments")

	cgr.POST("", ch.Create(), r.Authenticate, r.RBAC.HasRole(types.Registered))
	cgr.GET("/:pid", ch.ReadsByPostID())
	cgr.DELETE("/:id", ch.Delete(), r.Authenticate, r.RBAC.HasRole(types.Registered))
}
