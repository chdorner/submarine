package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/handler"
	"github.com/chdorner/submarine/router"
	"github.com/chdorner/submarine/test"
	"github.com/stretchr/testify/require"
)

func TestTagHandler(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	var publicTitles []string
	var privateTitles []string
	for i := 0; i < 8; i++ {
		title := fmt.Sprintf("Bookmark %d", i)

		var public bool
		if i%2 == 0 {
			public = true
			title = fmt.Sprintf("%s public", title)
			publicTitles = append(publicTitles, title)
		} else {
			public = false
			title = fmt.Sprintf("%s private", title)
			privateTitles = append(privateTitles, title)
		}
		_, err := repo.Create(data.BookmarkForm{
			URL:    fmt.Sprintf("https://example-%d.com", i),
			Title:  title,
			Public: public,
			Tags:   "articles",
		})
		require.NoError(t, err)
	}

	e := router.NewBaseApp(db)

	// queries public bookmarks when logged out
	req := httptest.NewRequest(http.MethodGet, "/tags/articles", strings.NewReader(""))
	rec := httptest.NewRecorder()
	sc := test.NewUnauthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("name")
	sc.SetParamValues("articles")

	err := handler.TagHandler(sc)
	require.NoError(t, err)
	for _, title := range publicTitles {
		require.Contains(t, rec.Body.String(), title)
	}
	for _, title := range privateTitles {
		require.NotContains(t, rec.Body.String(), title)
	}

	// queries all bookmarks when logged in
	req = httptest.NewRequest(http.MethodGet, "/tags/articles", strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("name")
	sc.SetParamValues("articles")

	err = handler.TagHandler(sc)
	require.NoError(t, err)
	for _, title := range append(publicTitles, privateTitles...) {
		require.Contains(t, rec.Body.String(), title)
	}

	// tag not found
	req = httptest.NewRequest(http.MethodGet, "/tags/missing", strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("name")
	sc.SetParamValues("missing")

	err = handler.TagHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
}
