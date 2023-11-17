package user

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/echoutils"
)

type Handler interface {
	Create() echo.HandlerFunc
	Read() echo.HandlerFunc
	Reads() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
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

		err := h.service.Create(r)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, "OK")
	}
}

func (h *handler) Read() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := new(dto.RequestWithID)
		if err := echoutils.BindAndValidate(c, r); err != nil {
			return err
		}

		u, err := h.service.Read(r)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, u)
	}
}

func (h *handler) Reads() echo.HandlerFunc {
	return func(c echo.Context) error {
		p := pagination.NewPagination()
		if err := echoutils.BindAndValidate(c, p); err != nil {
			return err
		}

		res, err := h.service.Reads(p)
		if err != nil {
			return err
		}

		p.PaginationHeader(c)

		return c.JSON(http.StatusOK, res)
	}
}

func (h *handler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := new(dto.UserUpdateRequest)
		if err := echoutils.BindAndValidate(c, r); err != nil {
			return err
		}

		res, err := h.service.Update(c.Request().Context(), r)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *handler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := new(dto.RequestWithID)
		if err := echoutils.BindAndValidate(c, r); err != nil {
			return err
		}

		res, err := h.service.Delete(r)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	}
}
