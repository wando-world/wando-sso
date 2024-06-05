package router

import (
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/api"
)

func SetupEnvRoutes(g *echo.Group) {
	g.GET("/env", api.EnvHandler)
}
