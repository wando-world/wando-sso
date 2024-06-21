package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wando-world/wando-sso/internal/config"
	"github.com/wando-world/wando-sso/internal/db"
	"github.com/wando-world/wando-sso/internal/router"
	"github.com/wando-world/wando-sso/internal/utils"
)

func main() {
	cfg := config.New()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = utils.NewValidator()

	db.InitDB(cfg.DbUrl)
	jwtUtils := utils.NewJwtUtils(cfg.ATKSecret, cfg.RTKSecret)
	router.SetupRoutes(e, jwtUtils)

	e.Logger.Fatal(e.Start(cfg.Port))
}
