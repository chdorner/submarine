package handler

import (
	"net/http"

	"github.com/chdorner/submarine/middleware"
	"github.com/labstack/echo/v4"
)

func SettingsHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	return sc.Render(http.StatusOK, "settings.html", map[string]interface{}{})
}
