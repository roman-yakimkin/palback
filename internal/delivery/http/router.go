package http

import (
	"net/http"
	"palback/internal/usecase/port"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"palback/internal/config"
	mwApp "palback/internal/delivery/http/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

// Validate реализует интерфейс echo.Validator
func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewRouter(
	cfg *config.Config,
	authenticator Authenticator,
	rateLimiter port.RateLimiter,
	countryHandler *CountryHandler,
	regionHandler *RegionHandler,
	cityTypeHandler *CityTypeHandler,
	placeTypeHandler *PlaceTypeHandler,
	userHandler *UserHandler,
) *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(mwApp.AuthMiddleware(authenticator))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.FrontendOrigin},
		AllowCredentials: true,
	}))

	// Работа со странами
	e.GET("/countries/:id", countryHandler.Get)
	e.GET("/countries", countryHandler.GetAll)
	e.POST("/countries", countryHandler.Post)
	e.PUT("/countries/:id", countryHandler.Put)
	e.DELETE("/countries/:id", countryHandler.Delete)
	e.POST("/countries/order", countryHandler.Order)

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

	// Работа с пользователями
	e.POST("/users/register", userHandler.Register,
		mwApp.RateLimitByIP(rateLimiter, 5*100, 600, "register"))
	e.POST("/users/verify-email", userHandler.VerifyEmail)
	e.POST("/users/resend-verification", userHandler.ResendVerification,
		mwApp.RateLimitByIP(rateLimiter, 5*100, 60, "resend-verification"))
	e.POST("/users/login", userHandler.Login,
		mwApp.RateLimitByIP(rateLimiter, 5*100, 60, "login"))
	e.POST("/users/logout", userHandler.Logout)
	e.POST("/users/reset-password", userHandler.ResetPassword,
		mwApp.RateLimitByIP(rateLimiter, 6*100, 3600, "reset"))
	e.POST("/users/reset-password/confirm", userHandler.ResetPasswordConfirm)
	e.GET("/users/me", userHandler.Me)
	e.GET("/users/profile", userHandler.ResetPassword)
	e.DELETE("users/delete", userHandler.Delete)

	return e
}
