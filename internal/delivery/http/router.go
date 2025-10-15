package http

import (
	"github.com/labstack/echo/v4"
)

func NewRouter(
	countryHandler *CountryHandler,
	regionHandler *RegionHandler,
	cityTypeHandler *CityTypeHandler,
	placeTypeHandler *PlaceTypeHandler,
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
	e.GET("/countries/:id/regions", regionHandler.GetByCountry)
	e.POST("/regions", regionHandler.Post)
	e.PUT("/regions/:id", regionHandler.Put)
	e.DELETE("/regions/:id", regionHandler.Delete)

	// Работа с типами населенных пунктов
	e.GET("/city-types/:id", cityTypeHandler.Get)
	e.GET("/city-types", cityTypeHandler.GetAll)

	// Работа с типами святых мест
	e.GET("/place-types/:id", placeTypeHandler.Get)
	e.GET("/place-types", placeTypeHandler.GetAll)

	return e
}
