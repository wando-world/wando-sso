package test

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/wando-world/wando-sso/internal/router"
	"testing"
)

func TestSetupRoutes(t *testing.T) {
	e := echo.New()
	router.SetupRoutes(e)

	require.NotNil(t, e.Routes())
}
