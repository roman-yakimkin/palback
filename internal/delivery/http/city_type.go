package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"palback/internal/delivery/http/dto"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
	"palback/internal/pkg/helpers"
	"palback/internal/usecase"
)

type CityTypeHandler struct {
	service usecase.CityTypeService
}

func NewCityTypeHandler(service usecase.CityTypeService) *CityTypeHandler {
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

	return c.JSON(http.StatusOK, dto.CreateCityTypeResponse(helpers.FromPtr(data)))
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

	return c.JSON(http.StatusOK, dto.CreateCityTypeResponseList(data))
}
