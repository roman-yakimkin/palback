package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"palback/internal/domain"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
)

type CountryHandler struct {
	service domain.CountryService
}

type country struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewCountryHandler(service domain.CountryService) *CountryHandler {
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

	return c.JSON(http.StatusOK, data)
}

// GetAll Получить все страны
func (h *CountryHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	data, err := h.service.GetAll(ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, data)
}

// Post Добавить страну
func (h *CountryHandler) Post(c echo.Context) error {
	ctx := c.Request().Context()

	var req country

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
		ID:   req.ID,
		Name: req.Name,
	})

	if err != nil {
		switch {
		case localErrors.IsOneOf(err, domain.ErrCountryAlreadyAdded, domain.ErrCountryNameNotUnique):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf(err.Error()),
			)
		}
	}

	return c.JSON(http.StatusOK, data)
}

// Put Изменить страну
func (h *CountryHandler) Put(c echo.Context) error {
	ctx := c.Request().Context()

	var req country

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
		Name: req.Name,
	})

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCountryAlreadyAdded):
			return echo.NewHTTPError(http.StatusConflict, "страна с данным id уже существует")
		case errors.Is(err, domain.ErrCountryNotFound):
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
