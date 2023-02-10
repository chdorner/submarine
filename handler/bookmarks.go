package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/middleware"
)

func BookmarksListHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	repo := data.NewBookmarkRepository(sc.DB)

	privacy := data.BookmarkPrivacyPublic
	if sc.IsAuthenticated() {
		privacy = data.BookmarkPrivacyQueryAll
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	result, err := repo.List(data.BookmarkListRequest{
		Privacy: privacy,
		Order:   "created_at desc",
		Offset:  offset,
	})
	if err != nil {
		return sc.Render(http.StatusOK, "bookmarks_list.html", map[string]interface{}{
			"error": "Failed to fetch bookmarks",
		})
	}

	err = sc.Render(http.StatusOK, "bookmarks_list.html", map[string]interface{}{
		"result": result,
	})

	return err
}

func BookmarksNewHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	if !sc.IsAuthenticated() {
		return sc.RedirectToLogin()
	}

	return sc.Render(http.StatusOK, "bookmarks_new.html", map[string]interface{}{
		"bookmark": data.BookmarkCreate{
			URL:         c.QueryParam("url"),
			Title:       c.QueryParam("title"),
			Description: c.QueryParam("desc"),
		},
	})
}

func BookmarksCreateHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	if !sc.IsAuthenticated() {
		return sc.RedirectToLogin()
	}

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
