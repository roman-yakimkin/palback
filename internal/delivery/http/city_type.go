package http

import (
	"errors"
	"net/http"
	"palback/internal/domain"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"

	"github.com/labstack/echo/v4"
)

type CityTypeHandler struct {
	service domain.CityTypeService
}

func NewCityTypeHandler(service domain.CityTypeService) *CityTypeHandler {
	return &CityTypeHandler{
		service: service,
	}
}

func (h *CityTypeHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	var data *model.CityType

	id, err := getPositiveIntParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка получения типа населенного пункта по id: "+err.Error())
	}

	data, err = h.service.Get(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, localErrors.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, data)
}

func (h *CityTypeHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	data, err := h.service.GetAll(ctx)

	if err != nil {
		switch {
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, data)
}
