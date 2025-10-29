package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetterUserID interface {
	GetUserID(context.Context, *http.Request) (int, error)
}

func AuthMiddleware(authenticator GetterUserID) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, _ := authenticator.GetUserID(c.Request().Context(), c.Request())
			if userID > 0 {
				c.Set("user_id", userID)
			}
			return next(c)
		}
	}
}
