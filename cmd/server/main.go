package main

import (
	"github.com/labstack/echo/v4"
	"github.com/wando-world/wando-sso/internal/config"
	"github.com/wando-world/wando-sso/internal/db"
	"github.com/wando-world/wando-sso/internal/router"
)

func main() {
	e := echo.New()

	cfg := config.New()
	db.InitDB(cfg.DbUrl)
	router.SetupRoutes(e)

	e.Logger.Fatal(e.Start(cfg.Port))
}
