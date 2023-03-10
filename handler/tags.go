package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/middleware"
	"github.com/labstack/echo/v4"
)

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
		Order:   "created_at desc",
		Offset:  offset,

		PaginationPathPrefix: fmt.Sprintf("/tags/%s?", tag.Name),
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
