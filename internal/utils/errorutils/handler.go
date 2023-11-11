package errorutils

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler custom HTTP error handler for echo framework.
func Handler(err error, c echo.Context) {
	var ae *APIError

	ok := errors.As(err, &ae)
	if ok {
		err = c.JSON(StatusCode(ae.Code), ae)
		if err != nil {
			c.Echo().Logger.Error(err)
		}

		return
	}

	var aes *APIErrors

	ok = errors.As(err, &aes)
	if ok {
		err = c.JSON(http.StatusBadRequest, aes)
		if err != nil {
			c.Echo().Logger.Error(err)
		}

		return
	}

	code := Code(err)

	err = c.JSON(StatusCode(code), &APIError{
		Code:    code,
		Message: err.Error(),
		Err:     err,
	})
	if err != nil {
		c.Echo().Logger.Error(err)
	}
}
