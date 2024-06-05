package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func EnvHandler(c echo.Context) error {
	env := os.Getenv("ENV")
	if env == "" {
		env = "test"
	}
	return c.String(http.StatusOK, env)
}
