package handler

import (
	"net/http"

	"github.com/chdorner/submarine/middleware"
	"github.com/labstack/echo/v4"
)

func SettingsHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	if !sc.IsAuthenticated() {
		return sc.RedirectToLogin()
	}

	return sc.Render(http.StatusOK, "settings.html", map[string]interface{}{
		"scheme": c.Scheme(),
		"host":   c.Request().Host,
	})
}
