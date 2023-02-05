package handler

import (
	_ "embed"
	"net/http"

	"github.com/chdorner/submarine/middleware"
	"github.com/labstack/echo/v4"
)

func LoginViewHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)

	return sc.Render(http.StatusOK, "login.html", map[string]interface{}{})
}
