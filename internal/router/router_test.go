package router

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetupRoutes(t *testing.T) {
	e := echo.New()
	SetupRoutes(e)

	require.NotNil(t, e.Routes())
}
