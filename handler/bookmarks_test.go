package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/handler"
	"github.com/chdorner/submarine/router"
	"github.com/chdorner/submarine/test"
	"github.com/stretchr/testify/require"
)

func TestBookmarksNewHandler(t *testing.T) {
	// unauthenticated
	e := router.NewBaseApp(nil)
	req := httptest.NewRequest(http.MethodGet, "/login", strings.NewReader(""))
	rec := httptest.NewRecorder()
	sc := test.NewUnauthenticatedContext(e.NewContext(req, rec), nil)
	err := handler.BookmarksNewHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.True(t, strings.HasPrefix(rec.Result().Header.Get("Location"), "/login?next="))
}

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
	sc := test.NewAuthenticatedContext(e.NewContext(req, rec), db)

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
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)

	err = handler.BookmarksCreateHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Contains(t, rec.Body.String(), "Failed to create bookmark")
	require.Contains(t, rec.Body.String(), "URL is required")

	// unauthenticated
	form = url.Values{}
	form.Add("url", "https://example.com/about")
	req = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)
	rec = httptest.NewRecorder()
	sc = test.NewUnauthenticatedContext(e.NewContext(req, rec), db)
	err = handler.BookmarksCreateHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.True(t, strings.HasPrefix(rec.Result().Header.Get("Location"), "/login"))
}

func TestBookmarksListHandler(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	var publicTitles []string
	var privateTitles []string
	for i := 0; i < 8; i++ {
		title := fmt.Sprintf("Bookmark %d", i)

		var privacy data.BookmarkPrivacy
		if i%2 == 0 {
			privacy = data.BookmarkPrivacyPublic
			title = fmt.Sprintf("%s public", title)
			publicTitles = append(publicTitles, title)
		} else {
			privacy = data.BookmarkPrivacyPrivate
			title = fmt.Sprintf("%s private", title)
			privateTitles = append(privateTitles, title)
		}
		_, err := repo.Create(data.BookmarkCreate{
			URL:     fmt.Sprintf("https://example-%d.com", i),
			Title:   title,
			Privacy: privacy,
		})
		require.NoError(t, err)
	}

	e := router.NewBaseApp(db)

	// queries public bookmarks when logged out
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	sc := test.NewUnauthenticatedContext(e.NewContext(req, rec), db)

	err := handler.BookmarksListHandler(sc)
	require.NoError(t, err)
	for _, title := range publicTitles {
		require.Contains(t, rec.Body.String(), title)
	}
	for _, title := range privateTitles {
		require.NotContains(t, rec.Body.String(), title)
	}

	// queries all bookmarks when logged in
	req = httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)

	err = handler.BookmarksListHandler(sc)
	require.NoError(t, err)
	for _, title := range append(publicTitles, privateTitles...) {
		require.Contains(t, rec.Body.String(), title)
	}
}
