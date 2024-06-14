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
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = utils.NewValidator()

	cfg := config.New()
	db.InitDB(cfg.DbUrl)
	router.SetupRoutes(e)

	e.Logger.Fatal(e.Start(cfg.Port))
}
