package user

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/echoutils"
)

type Handler interface {
	Create() echo.HandlerFunc
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := new(dto.UserCreateRequest)
		if err := echoutils.BindAndValidate(c, r); err != nil {
			return err
		}

		ctx := c.Request().Context()

		//claim, err := linkifycontext.LinkifyUser(ctx)
		//if err != nil {
		//	return err
		//}

		u, err := h.service.Create(ctx, r)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, u)
	}
}
