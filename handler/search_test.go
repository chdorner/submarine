package handler_test

import (
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

func TestSearchHandler(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	_, err := repo.Create(data.BookmarkForm{
		URL:   "https://en.wikipedia.org",
		Title: "Wikipedia",
		Tags:  "toRead",
	})
	require.NoError(t, err)
	_, err = repo.Create(data.BookmarkForm{
		URL:   "https://example.com",
		Title: "Example",
	})
	require.NoError(t, err)

	e := router.NewBaseApp(db)

	// tag search
	q := make(url.Values)
	q.Set("q", "toread")
	req := httptest.NewRequest(http.MethodGet, "/search?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	sc := test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	err = handler.SearchHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Contains(t, rec.Body.String(), "Matching Tags")
	require.Contains(t, rec.Body.String(), "toRead")

	// unauthenticated
	req = httptest.NewRequest(http.MethodGet, "/search", nil)
	rec = httptest.NewRecorder()
	sc = test.NewUnauthenticatedContext(e.NewContext(req, rec), db)
	err = handler.SearchHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.True(t, strings.HasPrefix(rec.Result().Header.Get("Location"), "/login"))
}
