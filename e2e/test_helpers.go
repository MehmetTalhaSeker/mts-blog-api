package e2e

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/appcontext"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/validatorutils"
)

const baseURL = "http://localhost:8080/v1"

func InitEcho(middlewares ...func(next echo.HandlerFunc) echo.HandlerFunc) *echo.Echo {
	e := echo.New()
	e.Validator = validatorutils.NewValidator()
	e.HTTPErrorHandler = errorutils.Handler

	for _, middleware := range middlewares {
		e.Use(middleware)
	}

	return e
}

func AuthMid() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

func AuthMidUser(e *echo.Echo, u *model.User) {
	mid := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := &dto.Claims{
				UID:      u.ID,
				Role:     u.Role,
				Username: u.Username,
				Email:    u.Email,
			}

			// Use custom context functions to store values
			ctx := appcontext.WithMtsBlogUser(c.Request().Context(), claims)
			ctx = appcontext.WithMtsBlogRole(ctx, u.Role)
			ctx = appcontext.WithMtsBlogUserID(ctx, u.ID)

			// Update the request context
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
	e.Use(mid)
}

func ClearAuthMidUser(e *echo.Echo) {
	mid := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Update the request context
			c.SetRequest(c.Request().WithContext(context.Background()))

			return next(c)
		}
	}
	e.Use(mid)
}

func Get(ctx context.Context, path string, headers ...map[string]string) (int, []byte, http.Header, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+path, nil)
	if err != nil {
		return 0, nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	for _, h := range headers {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, err
	}

	return resp.StatusCode, body, resp.Header, nil
}

func Post(ctx context.Context, path string, data []byte, headers ...map[string]string) (int, []byte, http.Header, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+path, bytes.NewBuffer(data))
	if err != nil {
		return 0, nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	for _, h := range headers {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, err
	}

	return resp.StatusCode, body, resp.Header, nil
}

func Put(ctx context.Context, path string, data []byte, headers ...map[string]string) (int, []byte, http.Header, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, baseURL+path, bytes.NewBuffer(data))
	if err != nil {
		return 0, nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	for _, h := range headers {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, err
	}

	return resp.StatusCode, body, resp.Header, nil
}

func Delete(ctx context.Context, path string, headers ...map[string]string) (int, []byte, http.Header, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, baseURL+path, nil)
	if err != nil {
		return 0, nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	for _, h := range headers {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, err
	}

	return resp.StatusCode, body, resp.Header, nil
}

func ReverseAndSlice(arr []any, start, end int) []any {
	rev := reverse(arr)

	return GetSubSlice(rev, start, end)
}

func GetSubSlice(arr []any, start, end int) []any {
	if start < 0 {
		start = 0
	}

	if end > len(arr) {
		end = len(arr)
	}

	return arr[start:end]
}

func reverse(arr []any) []any {
	length := len(arr)

	result := make([]any, length)
	for i, v := range arr {
		result[length-i-1] = v
	}

	return result
}
