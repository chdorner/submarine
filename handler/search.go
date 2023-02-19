package handler

import (
	"net/http"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/middleware"
	"github.com/labstack/echo/v4"
)

func SearchHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	if !sc.IsAuthenticated() {
		return sc.RedirectToLogin()
	}

	tagRepo := data.NewTagRepository(sc.DB)

	var tags []data.Tag
	var bookmarks []data.Bookmark

	query := sc.QueryParam("q")
	if query != "" {
		tags, _ = tagRepo.Search(query)
	}

	return sc.Render(http.StatusOK, "search.html", map[string]interface{}{
		"query":     query,
		"tags":      tags,
		"bookmarks": bookmarks,
	})
}
