package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func NewHealthzHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}
}
