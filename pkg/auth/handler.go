package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/echoutils"
)

type Handler interface {
	Login() echo.HandlerFunc
	Register() echo.HandlerFunc
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := new(dto.LoginRequest)
		if err := echoutils.BindAndValidate(c, r); err != nil {
			return err
		}

		resp, err := h.service.Login(r)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func (h *handler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := new(dto.RegisterRequest)
		if err := echoutils.BindAndValidate(c, r); err != nil {
			return err
		}

		resp, err := h.service.Register(r)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, resp)
	}
}
