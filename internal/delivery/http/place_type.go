package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"palback/internal/app"
	"palback/internal/delivery/http/dto"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
	"palback/internal/pkg/helpers"
)

type PlaceTypeHandler struct {
	service app.PlaceTypeService
}

func NewPlaceTypeHandler(service app.PlaceTypeService) *PlaceTypeHandler {
	return &PlaceTypeHandler{
		service: service,
	}
}

func (h *PlaceTypeHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	var data *model.PlaceType

	id, err := getPositiveIntParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка получения типа святого места по id: "+err.Error())
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

	return c.JSON(http.StatusOK, dto.CreatePlaceTypeResponse(helpers.FromPtr(data)))
}

func (h *PlaceTypeHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	data, err := h.service.GetAll(ctx)

	if err != nil {
		switch {
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, dto.CreatePlaceTypeResponseList(data))
}
