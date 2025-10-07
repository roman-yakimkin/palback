package http

import (
	"github.com/labstack/echo/v4"
)

func NewRouter(
	countryHandler *CountryHandler,
	regionHandler *RegionHandler,
) *echo.Echo {
	e := echo.New()

	// Работа со странами
	e.GET("/countries/:id", countryHandler.Get)
	e.GET("/countries", countryHandler.GetAll)
	e.POST("/countries", countryHandler.Post)
	e.PUT("/countries/:id", countryHandler.Put)
	e.DELETE("/countries/:id", countryHandler.Delete)

	// Работа с регионами
	e.GET("/regions/:id", regionHandler.Get)
	e.GET("/regions/by-country/:countryId", regionHandler.GetByCountry)
	e.POST("/regions", regionHandler.Post)
	e.PUT("/regions/:id", regionHandler.Put)
	e.DELETE("/regions/:id", regionHandler.Delete)

	return e
}
