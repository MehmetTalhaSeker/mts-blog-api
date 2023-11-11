package echoutils

import (
	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

func BindAndValidate(c echo.Context, i any) error {
	if err := c.Bind(i); err != nil {
		return errorutils.New(errorutils.ErrBinding, err)
	}

	return c.Validate(i)
}
