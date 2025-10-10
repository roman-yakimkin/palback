package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"palback/internal/app"
	appModel "palback/internal/app/model"
	"palback/internal/delivery/http/dto"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
	"palback/internal/pkg/helpers"
)

type RegionHandler struct {
	service app.RegionService
}

func NewRegionHandler(service app.RegionService) *RegionHandler {
	return &RegionHandler{
		service: service,
	}
}

func (h *RegionHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	var data *appModel.RegionDetail

	id, err := getPositiveIntParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка получения региона по id: "+err.Error())
	}

	data, err = h.service.Get(ctx, id)

	if errors.Is(err, localErrors.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dto.CreateRegionResponse(helpers.FromPtr(data)))
}

func (h *RegionHandler) GetByCountry(c echo.Context) error {
	ctx := c.Request().Context()

	data, err := h.service.GetByCountry(ctx, c.Param("countryId"))

	if err != nil {
		switch {
		case errors.Is(err, app.ErrCountryNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, dto.CreateRegionResponseList(data))
}

func (h *RegionHandler) Post(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.RegionPostRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	data, err := h.service.Create(ctx, model.Region{
		CountryID: req.CountryID,
		Name:      req.Name,
	})

	if err != nil {
		switch {
		case localErrors.IsOneOf(err, app.ErrCountryHasNotRegions, app.ErrRegionNotUnique):
			return echo.NewHTTPError(
				http.StatusBadRequest,
				fmt.Sprintf(err.Error()),
			)
		default:
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("невозможно добавить регион: %s", err.Error()),
			)
		}
	}

	return c.JSON(http.StatusOK, dto.CreateRegionResponse(helpers.FromPtr(data)))
}

func (h *RegionHandler) Put(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := getPositiveIntParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка получения региона по id: "+err.Error())
	}

	var req dto.RegionPutRequest

	if err = c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.service.Update(ctx, id, model.Region{
		CountryID: req.CountryID,
		Name:      req.Name,
	})

	if err != nil {
		switch {
		case errors.Is(err, app.ErrCountryNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case localErrors.IsOneOf(err, app.ErrCountryHasNotRegions, app.ErrRegionNotUnique):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("невозможно изменить регион: %s", err.Error()),
			)
		}
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "регион обновлен"})
}

func (h *RegionHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := getPositiveIntParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка получения региона по id: "+err.Error())
	}

	err = h.service.Delete(ctx, id)

	if err != nil {
		switch {
		case localErrors.IsOneOf(err, app.ErrRegionNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("невозможно удалить регион: %s", err.Error()),
			)
		}
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "регион удален"})
}
