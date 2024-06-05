package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func HomeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "SSO Server is running")
}
