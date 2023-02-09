package handler_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/handler"
	"github.com/chdorner/submarine/middleware"
	"github.com/chdorner/submarine/router"
	"github.com/chdorner/submarine/test"
	"github.com/stretchr/testify/require"
)

func TestBookmarksCreateHandler(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()

	contentType := "application/x-www-form-urlencoded"

	e := router.NewBaseApp(db)

	// success
	form := url.Values{}
	form.Add("url", "https://example.com/about")
	form.Add("title", "Example - About")
	form.Add("description", "About example.com")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()
	sc := middleware.InitSubmarineContext(e.NewContext(req, rec), db)

	err := handler.BookmarksCreateHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.Equal(t, "/", rec.Header().Get("Location"))

	var bookmark data.Bookmark
	result := db.Last(&bookmark)
	require.NoError(t, result.Error)
	require.Equal(t, form.Get("url"), bookmark.URL)

	// validation error
	form = url.Values{}
	form.Add("title", "Example - About")

	req = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)
	rec = httptest.NewRecorder()
	sc = middleware.InitSubmarineContext(e.NewContext(req, rec), db)

	err = handler.BookmarksCreateHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Contains(t, rec.Body.String(), "Failed to create bookmark")
	require.Contains(t, rec.Body.String(), "URL is required")

	// TODO: test unauthenticated access
}
