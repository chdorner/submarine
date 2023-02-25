package handler

import (
	"net/http"
	"net/url"
	"strconv"

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
	bookmarkRepo := data.NewBookmarkRepository(sc.DB)

	var tags []data.Tag
	var bookmarkResults *data.BookmarkSearchResponse

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	query := sc.QueryParam("q")

	params := url.Values{}
	params.Add("q", query)

	if query != "" {
		tags, _ = tagRepo.Search(query)
		bookmarkResults, _ = bookmarkRepo.Search(data.BookmarkSearchRequest{
			Query:                query,
			Offset:               offset,
			PaginationPathPrefix: "/search?" + params.Encode() + "&",
		})
	}

	return sc.Render(http.StatusOK, "search.html", map[string]interface{}{
		"query":      query,
		"tags":       tags,
		"hasResults": bookmarkResults != nil && bookmarkResults.Count > int64(0),
		"result":     bookmarkResults,
	})
}
