package router

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wando-world/wando-sso/internal/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnvRouter(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/sso/api/v1/env", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, api.EnvHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		expected := "test"
		require.Contains(t, rec.Body.String(), expected)
	}
}
