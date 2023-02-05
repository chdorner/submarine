package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func root(c echo.Context) error {
	sc := c.(*SubmarineContext)

	return sc.String(http.StatusOK, "Hello from submarine")
}
