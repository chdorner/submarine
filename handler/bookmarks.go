package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/chdorner/submarine/middleware"
)

func BookmarksListHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)

	err := sc.Render(http.StatusOK, "bookmarks_list.html", map[string]interface{}{})

	return err
}
