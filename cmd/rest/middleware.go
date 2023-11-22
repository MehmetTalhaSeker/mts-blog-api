package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	postgresadapter "github.com/MehmetTalhaSeker/mts-blog-api/internal/adapter/postgres"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/appcontext"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

func (app *application) authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errorutils.New(errorutils.ErrMissingAuthHeader, errorutils.ErrMissingAuthHeader)
			}

			ts := strings.Replace(authHeader, "Bearer ", "", 1)

			token, err := validateJWT(ts)
			if err != nil {
				return errorutils.New(errorutils.ErrInvalidToken, err)
			}

			mc, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return errorutils.New(errorutils.ErrUnauthorized, err)
			}

			dts, ok := mc["expiresAt"].(string)
			if !ok {
				return errorutils.New(errorutils.ErrInvalidToken, err)
			}

			exp, err := time.Parse(time.RFC3339, dts)
			if err != nil {
				return errorutils.New(errorutils.ErrInvalidToken, err)
			}

			if exp.Unix() < time.Now().Unix() {
				return errorutils.New(errorutils.ErrExpiredToken, err)
			}

			claims := &dto.Claims{}

			err = apputils.InterfaceToStruct(mc, claims)
			if err != nil {
				return err
			}

			ur := postgresadapter.NewUserRepository(app.db)

			u, err := ur.Read(claims.UID)
			if err != nil {
				return err
			}

			if u.Role != claims.Role || u.Status == types.Passive {
				return errorutils.New(errorutils.ErrInvalidRequest, nil)
			}

			// Use custom context functions to store values
			ctx := appcontext.WithMtsBlogUser(c.Request().Context(), claims)
			ctx = appcontext.WithMtsBlogRole(ctx, claims.Role)
			ctx = appcontext.WithMtsBlogUserID(ctx, claims.UID)

			// Update the request context
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sign- in method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
