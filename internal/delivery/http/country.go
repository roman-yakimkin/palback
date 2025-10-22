package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"palback/internal/app"
	"palback/internal/delivery/http/dto"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
	"palback/internal/pkg/helpers"
)

type CountryHandler struct {
	service app.CountryService
}

func NewCountryHandler(service app.CountryService) *CountryHandler {
	return &CountryHandler{
		service: service,
	}
}

// Get Получить одну страну по id
func (h *CountryHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	var data *model.Country

	data, err := h.service.Get(ctx, c.Param("id"))

	if errors.Is(err, localErrors.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, dto.CreateCountryResponse(helpers.FromPtr(data)))
}

// GetAll Получить все страны
func (h *CountryHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	data, err := h.service.GetAll(ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, dto.CreateCountryResponseList(data))
}

// Post Добавить страну
func (h *CountryHandler) Post(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.CountryPostRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !validCountryID.MatchString(req.ID) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"id страны должен состоять из строчных латинских букв и иметь в длину от 2 до 6 символов включительно",
		)
	}

	data, err := h.service.Create(ctx, model.Country{
		ID:         req.ID,
		Name:       req.Name,
		HasRegions: req.HasRegions,
		Weight:     req.Weight,
	})

	if err != nil {
		switch {
		case localErrors.IsOneOf(err, app.ErrCountryAlreadyAdded, app.ErrCountryNameNotUnique):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf(err.Error()),
			)
		}
	}

	dataRec := helpers.FromPtr(data)

	c.Response().Header().Set("location", "/countries/"+dataRec.ID)

	return c.JSON(http.StatusCreated, dto.CreateCountryResponse(dataRec))
}

// Put Изменить страну
func (h *CountryHandler) Put(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.CountryPutRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id := c.Param("id")
	if !validCountryID.MatchString(id) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"id страны должен состоять из строчных латинских букв и иметь в длину от 2 до 6 символов включительно",
		)
	}

	err := h.service.Update(ctx, id, model.Country{
		Name:       req.Name,
		HasRegions: req.HasRegions,
		Weight:     req.Weight,
	})

	if err != nil {
		switch {
		case errors.Is(err, app.ErrCountryAlreadyAdded):
			return echo.NewHTTPError(http.StatusConflict, "страна с данным id уже существует")
		case errors.Is(err, app.ErrCountryNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("невозможно изменить страну: %s", err.Error()),
			)
		}
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "страна обновлена"})
}

// Delete Удалить страну
func (h *CountryHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	err := h.service.Delete(ctx, c.Param("id"))

	if errors.Is(err, localErrors.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("невозможно удалить страну: %s", err.Error()),
		)
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "страна удалена"})

}

func (h *CountryHandler) Order(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.CountryOrderRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.service.Order(ctx, req.Order)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("невозможно упорядочить страны: %s", err.Error()),
		)
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "страны упорядочены"})
}
