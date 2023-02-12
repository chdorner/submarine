package handler

import (
	"net/http"
	"strconv"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/middleware"
	"github.com/labstack/echo/v4"
)

func TagsListHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	return sc.Render(http.StatusOK, "tags_list.html", map[string]interface{}{})
}

func TagHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	repo := data.NewTagRepository(sc.DB)

	tag, err := repo.GetByName(sc.Param("name"))
	if err != nil || tag == nil {
		return sc.RenderNotFound()
	}

	privacy := data.BookmarkPrivacyPublic
	if sc.IsAuthenticated() {
		privacy = data.BookmarkPrivacyQueryAll
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	bookmarksRepo := data.NewBookmarkRepository(sc.DB)
	result, err := bookmarksRepo.List(data.BookmarkListRequest{
		Privacy: privacy,
		TagID:   tag.ID,
		Offset:  offset,
	})
	if err != nil {
		return sc.Render(http.StatusOK, "tags_show.html", map[string]interface{}{
			"error": "Failed to fetch bookmarks",
		})
	}

	return sc.Render(http.StatusOK, "tags_show.html", map[string]interface{}{
		"tag":    tag,
		"result": result,
	})
}
