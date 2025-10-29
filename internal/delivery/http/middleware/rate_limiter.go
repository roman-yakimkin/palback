package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"palback/internal/usecase/port"
)

func RateLimitByIP(
	limiter port.RateLimiter,
	limit int,
	windowSeconds int,
	keyPrefix string,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			key := fmt.Sprintf("rl:%s:%s", keyPrefix, ip)

			allowed, err := limiter.Allow(c.Request().Context(), key, limit, windowSeconds)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "rate limit check failed")
			}
			if !allowed {
				return echo.NewHTTPError(http.StatusTooManyRequests, "cлишком много запросов с одного IP-адреса")
			}

			return next(c)
		}
	}
}
