package auth

import (
	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/repository"
)

type Router struct {
	RouterGroup    *echo.Group
	UserRepository repository.User
}

func (r *Router) New() {
	as := NewService(r.UserRepository)
	ah := NewHandler(as)

	ugr := r.RouterGroup.Group("/auth")

	ugr.POST("/login", ah.Login())
	ugr.POST("/register", ah.Register())
}
