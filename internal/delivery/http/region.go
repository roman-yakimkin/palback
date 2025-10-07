package http

import (
	"errors"
	"fmt"
	"net/http"
	"palback/internal/domain"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"

	"github.com/labstack/echo/v4"
)

type RegionHandler struct {
	service domain.RegionService
}

func NewRegionHandler(service domain.RegionService) *RegionHandler {
	return &RegionHandler{
		service: service,
	}
}

type region struct {
	CountryID string `json:"country_id"`
	Name      string `json:"name"`
}

func (h *RegionHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	var data *model.Region

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

	return c.JSON(http.StatusOK, data)
}

func (h *RegionHandler) GetByCountry(c echo.Context) error {
	ctx := c.Request().Context()

	data, err := h.service.GetByCountry(ctx, c.Param("countryId"))

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCountryNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, data)
}

func (h *RegionHandler) Post(c echo.Context) error {
	ctx := c.Request().Context()

	var req region

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	data, err := h.service.Create(ctx, model.Region{
		CountryID: req.CountryID,
		Name:      req.Name,
	})

	if err != nil {
		switch {
		case localErrors.IsOneOf(err, domain.ErrCountryHasNotRegions, domain.ErrRegionNotUnique):
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

	return c.JSON(http.StatusOK, data)
}

func (h *RegionHandler) Put(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := getPositiveIntParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка получения региона по id: "+err.Error())
	}

	var req region

	if err = c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.service.Update(ctx, id, model.Region{
		CountryID: req.CountryID,
		Name:      req.Name,
	})

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCountryNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case localErrors.IsOneOf(err, domain.ErrCountryHasNotRegions, domain.ErrRegionNotUnique):
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
		case localErrors.IsOneOf(err, domain.ErrRegionNotFound):
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
