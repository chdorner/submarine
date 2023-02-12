package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

		PaginationPathPrefix: "/?",
	})
	if err != nil {
		return sc.Render(http.StatusOK, "bookmarks_list.html", map[string]interface{}{
			"error": "Failed to fetch bookmarks.",
		})
	}

	err = sc.Render(http.StatusOK, "bookmarks_list.html", map[string]interface{}{
		"result": result,
	})

	return err
}

func BookmarkShowHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	repo := data.NewBookmarkRepository(sc.DB)

	id, err := strconv.Atoi(sc.Param("id"))
	if err != nil {
		return sc.RenderNotFound()
	}
	bookmark, err := repo.Get(uint(id))
	if err != nil || bookmark == nil {
		return sc.RenderNotFound()
	}
	if bookmark.Privacy == data.BookmarkPrivacyPrivate && !sc.IsAuthenticated() {
		return sc.RenderNotFound()
	}

	return sc.Render(http.StatusOK, "bookmarks_show.html", map[string]interface{}{
		"bookmark": bookmark,
	})
}

func BookmarksNewHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	if !sc.IsAuthenticated() {
		return sc.RedirectToLogin()
	}

	return sc.Render(http.StatusOK, "bookmarks_new.html", map[string]interface{}{
		"form": map[string]string{
			"submit": "Create Bookmark",
			"action": "/bookmarks",
		},
		"bookmark": data.BookmarkForm{
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

	formTplData := map[string]string{
		"submit": "Create Bookmark",
		"action": "/bookmarks",
	}

	repo := data.NewBookmarkRepository(sc.DB)

	form, validationErr := parseAndValidateForm(sc)
	if validationErr != nil {
		return sc.Render(http.StatusOK, "bookmarks_new.html", map[string]interface{}{
			"form":             formTplData,
			"error":            "Failed to create bookmark",
			"validationErrors": validationErr.Fields,
			"bookmark":         form,
		})
	}

	bookmark, err := repo.Create(*form)
	if err != nil {
		return sc.Render(http.StatusOK, "bookmarks_new.html", map[string]interface{}{
			"form":     formTplData,
			"error":    "Unexpected error happened when creating bookmark, please try again.",
			"bookmark": form,
		})
	}

	return sc.Redirect(http.StatusFound, fmt.Sprintf("/bookmarks/%d", bookmark.ID))
}

func BookmarkDeleteHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	if !sc.IsAuthenticated() {
		return sc.RedirectToLogin()
	}

	repo := data.NewBookmarkRepository(sc.DB)
	id, err := strconv.Atoi(sc.Param("id"))
	if err != nil {
		return sc.RenderNotFound()
	}
	err = repo.Delete(uint(id))
	if err != nil {
		return sc.RenderNotFound()
	}

	return sc.Redirect(http.StatusFound, "/")
}

func BookmarkEditViewHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	if !sc.IsAuthenticated() {
		return sc.RedirectToLogin()
	}

	repo := data.NewBookmarkRepository(sc.DB)
	id, err := strconv.Atoi(sc.Param("id"))
	if err != nil {
		return sc.RenderNotFound()
	}
	bookmark, err := repo.Get(uint(id))
	if err != nil || bookmark == nil {
		return sc.RenderNotFound()
	}

	tagNames := []string{}
	for _, tag := range bookmark.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	form := data.BookmarkForm{
		URL:         bookmark.URL,
		Title:       bookmark.Title,
		Description: bookmark.Description,
		Public:      bookmark.IsPublic(),
		Tags:        strings.Join(tagNames, ", "),
	}

	return sc.Render(http.StatusOK, "bookmarks_edit.html", map[string]interface{}{
		"form": map[string]string{
			"submit": "Edit Bookmark",
			"action": fmt.Sprintf("/bookmarks/%d/edit", bookmark.ID),
		},
		"bookmark": form,
	})
}

func BookmarkEditHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)
	if !sc.IsAuthenticated() {
		return sc.RedirectToLogin()
	}

	repo := data.NewBookmarkRepository(sc.DB)
	id, err := strconv.Atoi(sc.Param("id"))
	if err != nil {
		return sc.RenderNotFound()
	}
	bookmark, err := repo.Get(uint(id))
	if err != nil || bookmark == nil {
		return sc.RenderNotFound()
	}

	formTplData := map[string]string{
		"submit": "Edit Bookmark",
		"action": fmt.Sprintf("/bookmarks/%d/edit", bookmark.ID),
	}

	form, validationErr := parseAndValidateForm(sc)
	if validationErr != nil {
		return sc.Render(http.StatusOK, "bookmarks_edit.html", map[string]interface{}{
			"form":             formTplData,
			"error":            "Failed to edit bookmark",
			"validationErrors": validationErr.Fields,
			"bookmark":         form,
		})
	}

	err = repo.Update(bookmark.ID, *form)
	if err != nil {
		return sc.Render(http.StatusOK, "bookmarks_edit.html", map[string]interface{}{
			"form":     formTplData,
			"error":    "Unexpected error happened when editing bookmark, please try again.",
			"bookmark": form,
		})
	}

	return sc.Redirect(http.StatusFound, fmt.Sprintf("/bookmarks/%d", bookmark.ID))
}

func parseAndValidateForm(sc *middleware.SubmarineContext) (*data.BookmarkForm, *data.ValidationError) {
	public := false
	if sc.FormValue("public") == "on" {
		public = true
	}
	req := data.BookmarkForm{
		URL:         sc.FormValue("url"),
		Title:       sc.FormValue("title"),
		Description: sc.FormValue("description"),
		Public:      public,
		Tags:        sc.FormValue("tags"),
	}
	validationErr := req.IsValid()
	if validationErr != nil {
		return &req, validationErr
	}
	return &req, nil
}
