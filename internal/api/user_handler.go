package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateUserHandler(c echo.Context) error {
	return c.JSON(http.StatusCreated, "개발예정")
}
