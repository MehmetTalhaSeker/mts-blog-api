package rbac

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/appcontext"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

type RBAC interface {
	HasRole(types.Role) func(next echo.HandlerFunc) echo.HandlerFunc
	CheckHasRole(userRole types.Role, requiredRole types.Role) bool
	IsAdminAuthorized(ctx context.Context) bool
	IsModAuthorized(context.Context) bool
	IsMe(context.Context, uint64) bool
	CheckRoleAndUser(context.Context, uint64, types.Role) (*dto.Claims, error)
}

var roleScores = map[types.Role]uint8{
	types.Admin:      15,
	types.Mod:        14,
	types.Registered: 13,
}

func New() RBAC {
	return &rbac{}
}

type rbac struct{}

func (r *rbac) IsMe(ctx context.Context, userID uint64) bool {
	claims, err := appcontext.MtsBlogUser(ctx)
	if err != nil {
		return false
	}

	return userID == claims.UID
}

func (r *rbac) IsAdminAuthorized(ctx context.Context) bool {
	claims, err := appcontext.MtsBlogUser(ctx)
	if err != nil {
		return false
	}

	return roleScores[claims.Role] >= roleScores[types.Admin]
}

func (r *rbac) IsModAuthorized(ctx context.Context) bool {
	claims, err := appcontext.MtsBlogUser(ctx)
	if err != nil {
		return false
	}

	return roleScores[claims.Role] >= roleScores[types.Mod]
}

func (r *rbac) HasRole(routeRole types.Role) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, err := appcontext.MtsBlogRole(c.Request().Context())
			if err != nil {
				return errorutils.New(errorutils.ErrUnauthorized, nil)
			}

			if roleScores[userRole] < roleScores[routeRole] {
				return errorutils.New(errorutils.ErrUnauthorized, nil)
			}

			return next(c)
		}
	}
}

func (r *rbac) CheckHasRole(userRole types.Role, requiredRole types.Role) bool {
	return roleScores[userRole] >= roleScores[requiredRole]
}

func (r *rbac) CheckRoleAndUser(ctx context.Context, userID uint64, requiredRole types.Role) (*dto.Claims, error) {
	claims, err := appcontext.MtsBlogUser(ctx)
	if err != nil {
		return nil, err
	}

	if userID != claims.UID && !r.CheckHasRole(claims.Role, requiredRole) {
		return nil, errorutils.New(errorutils.ErrUnauthorized, nil)
	}

	return claims, nil
}
