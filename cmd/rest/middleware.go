package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/appcontext"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
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

			claims := &dto.Claims{}

			err = apputils.InterfaceToStruct(mc, claims)
			if err != nil {
				return err
			}

			fmt.Println(claims)

			//acc, err := s.GetAccountByID(id)
			//if err != nil {
			//	return errorutils.New(errorutils.ErrUserNotFound, err)
			//}
			//
			//
			//uc, ok := idToken.Claims["lu"]
			//
			//claims := &dto.Claims{}
			//
			//if ok {
			//	err = apputils.InterfaceToStruct(lu, claims)
			//	if err != nil {
			//		return err
			//	}
			//} else {
			//	fu, err := app.auth.GetUser(c.Request().Context(), idToken.UID)
			//	if err != nil {
			//		return errorutils.New(errorutils.ErrInvalidToken, err)
			//	}
			//
			//	var u *model.User
			//
			//	u, err = app.auth.VerifyUser(c.Request().Context(), idToken.Subject)
			//	if err != nil {
			//		return errorutils.New(errorutils.ErrUserNotFound, err)
			//	}
			//
			//	claims.UID = u.ID
			//	claims.Role = u.Role
			//	claims.Username = u.Username
			//
			//	if claims.Role == types.Disabled {
			//		return errorutils.New(errorutils.ErrUserDisabled, err)
			//	}
			//
			//	if claims == nil || claims.LID == "" {
			//		return errorutils.New(errorutils.ErrUserNotFound, err)
			//	}
			//}
			//

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
