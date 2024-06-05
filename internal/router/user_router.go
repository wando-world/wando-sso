package router

import (
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/api"
	"github.com/wando-world/wando-sso/internal/api/models"
	"net/http"
)

func SetupUserRoutes(g *echo.Group) {
	g.POST("/user", func(c echo.Context) error {
		var req models.CreateUserRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "입력값을 확인해주세요.")
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// 로직
		return api.CreateUserHandler(c)
	})
}
