package router

import (
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/api"
	"github.com/wando-world/wando-sso/internal/db"
	"github.com/wando-world/wando-sso/internal/mappers"
	"github.com/wando-world/wando-sso/internal/utils"
)

func SetupAuthRoutes(g *echo.Group, jwtUtils utils.IJwt) {
	passwordUtils := utils.NewPasswordUtils()
	authMapper := mappers.NewAuthMapper()
	userRepository := db.NewUserRepository(db.DB)
	authHandler := api.NewAuthHandler(passwordUtils, jwtUtils, authMapper, userRepository)

	g.POST("/login", authHandler.Login)
}
