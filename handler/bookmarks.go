package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/middleware"
)

func BookmarksListHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)

	err := sc.Render(http.StatusOK, "bookmarks_list.html", map[string]interface{}{})

	return err
}

func BookmarksNewHandler(c echo.Context) error {
	// TODO: require authentication
	sc := c.(*middleware.SubmarineContext)
	return sc.Render(http.StatusOK, "bookmarks_new.html", nil)
}

func BookmarksCreateHandler(c echo.Context) error {
	// TODO: require authentication
	sc := c.(*middleware.SubmarineContext)
	repo := data.NewBookmarkRepository(sc.DB)

	req := data.BookmarkCreate{
		URL:         sc.FormValue("url"),
		Title:       sc.FormValue("title"),
		Description: sc.FormValue("description"),
	}
	validationErr := req.IsValid()
	if validationErr != nil {
		return sc.Render(http.StatusOK, "bookmarks_new.html", map[string]interface{}{
			"error":            "Failed to create bookmark",
			"validationErrors": validationErr.Fields,
			"bookmark":         req,
		})
	}

	_, err := repo.Create(req)
	if err != nil {
		return sc.Render(http.StatusOK, "bookmarks_new.html", map[string]interface{}{
			"error": "Unexpected error happened when creating bookmark, please try again.",
		})
	}

	return sc.Redirect(http.StatusFound, "/")
}
