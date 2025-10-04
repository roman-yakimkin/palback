package http

import (
	"github.com/labstack/echo/v4"
)

func NewRouter(countryHandler *CountryHandler) *echo.Echo {
	e := echo.New()

	// Работа со странами
	e.GET("/countries/:id", countryHandler.Get)
	e.GET("/countries", countryHandler.GetAll)
	e.POST("/countries", countryHandler.Post)
	e.PUT("/countries/:id", countryHandler.Put)
	e.DELETE("/countries/:id", countryHandler.Delete)

	return e
}
