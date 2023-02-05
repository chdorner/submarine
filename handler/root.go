package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/chdorner/submarine/middleware"
)

func RootHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)

	name := "stranger"
	if sc.IsAuthenticated {
		name = "friend"
	}

	err := sc.Render(http.StatusOK, "root.html", map[string]interface{}{
		"name": name,
	})

	return err
}
