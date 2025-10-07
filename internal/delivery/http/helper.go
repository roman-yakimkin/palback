package http

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

var validCountryID = regexp.MustCompile(`^[a-z]{2,6}$`)

func getPositiveIntParam(c echo.Context, paramName string) (int, error) {
	paramStr := c.Param(paramName)
	if strings.TrimSpace(paramStr) == "" {
		return 0, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Отсутствует параметр %q", paramName))
	}

	num, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Неверный %q: должен быть числовым", paramName),
		)
	}

	if num <= 0 {
		return 0, echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Неверныйй %q: должен быть положительным числом", paramName),
		)
	}

	return num, nil
}

func getLang(c echo.Context) string {
	lang, ok := c.Get("lang").(string)
	if !ok {
		return "ru"
	}

	return lang
}
