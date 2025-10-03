package middleware

import (
	"github.com/labstack/echo/v4"
	"palback/internal/config"
	"strings"
)

// SetupLanguage Установить язык для запросов
func SetupLanguage() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			lang := c.QueryParam("lang")
			if strings.TrimSpace(lang) == "" {
				lang = config.GetLang()
			}

			c.Set("lang", lang)

			return next(c)
		}
	}
}
