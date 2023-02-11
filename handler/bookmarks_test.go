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
	form.Add("url", "https://example.com/public")
	form.Add("title", "Example - About")
	form.Add("description", "About example.com")
	form.Add("public", "on")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()
	sc := test.NewAuthenticatedContext(e.NewContext(req, rec), db)

	err := handler.BookmarksCreateHandler(sc)
	require.NoError(t, err)

	var bookmark data.Bookmark
	result := db.Last(&bookmark)
	require.NoError(t, result.Error)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.Equal(t, fmt.Sprintf("/bookmarks/%d", bookmark.ID), rec.Header().Get("Location"))

	require.Equal(t, form.Get("url"), bookmark.URL)
	require.Equal(t, data.BookmarkPrivacyPublic, bookmark.Privacy)

	// private bookmark (missing public in form values)
	form = url.Values{}
	form.Add("url", "https://example.com/private")

	req = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)

	err = handler.BookmarksCreateHandler(sc)
	require.NoError(t, err)

	bookmark = data.Bookmark{}
	result = db.Last(&bookmark)
	require.NoError(t, result.Error)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.Equal(t, fmt.Sprintf("/bookmarks/%d", bookmark.ID), rec.Header().Get("Location"))

	require.Equal(t, form.Get("url"), bookmark.URL)
	require.Equal(t, data.BookmarkPrivacyPrivate, bookmark.Privacy)

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

func TestBookmarksShowHandler(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	publicBookmark, err := repo.Create(data.BookmarkForm{
		URL:    "https://example.com/public",
		Title:  "Example public",
		Public: true,
	})
	require.NoError(t, err)
	privateBookmark, err := repo.Create(data.BookmarkForm{
		URL:    "https://example.com/private",
		Title:  "Example private",
		Public: false,
	})
	require.NoError(t, err)

	e := router.NewBaseApp(db)

	// public bookmark - unauthenticated
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/bookmarks/%d", publicBookmark.ID), strings.NewReader(""))
	rec := httptest.NewRecorder()
	sc := test.NewUnauthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues(fmt.Sprint(publicBookmark.ID))

	err = handler.BookmarkShowHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Contains(t, rec.Body.String(), publicBookmark.Title)

	// private bookmark - authenticated
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/bookmarks/%d", privateBookmark.ID), strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues(fmt.Sprint(privateBookmark.ID))

	err = handler.BookmarkShowHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Contains(t, rec.Body.String(), privateBookmark.Title)

	// private bookmark - unauthenticated
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/bookmarks/%d", privateBookmark.ID), strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewUnauthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues(fmt.Sprint(privateBookmark.ID))

	err = handler.BookmarkShowHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)

	// non-existing id
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/bookmarks/%d", privateBookmark.ID), strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues("42")

	err = handler.BookmarkShowHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)

	// non-int id
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/bookmarks/%d", privateBookmark.ID), strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues("notaninteger")

	err = handler.BookmarkShowHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
}

func TestBookmarkDeleteHandler(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	bookmark, err := repo.Create(data.BookmarkForm{
		URL:   "https://example.com/public",
		Title: "Example public",
	})
	require.NoError(t, err)

	e := router.NewBaseApp(db)

	// delete unauthenticated
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/bookmarks/%d/delete", bookmark.ID), strings.NewReader(""))
	rec := httptest.NewRecorder()
	sc := test.NewUnauthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues(fmt.Sprint(bookmark.ID))

	err = handler.BookmarkDeleteHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.Equal(t, "/login", rec.Result().Header.Get("Location"))

	existing, err := repo.Get(bookmark.ID)
	require.NoError(t, err)
	require.NotNil(t, existing)

	// delete authenticated
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/bookmarks/%d/delete", bookmark.ID), strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues(fmt.Sprint(bookmark.ID))

	err = handler.BookmarkDeleteHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.Equal(t, "/", rec.Result().Header.Get("Location"))

	existing, err = repo.Get(bookmark.ID)
	require.NoError(t, err)
	require.Nil(t, existing)

	// delete non-existing bookmark
	req = httptest.NewRequest(http.MethodPost, "/bookmarks/42/delete", strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues("42")

	err = handler.BookmarkDeleteHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)

	// delete non-existing bookmark with non-integer id
	req = httptest.NewRequest(http.MethodPost, "/bookmarks/notaninteger/delete", strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues("notaninteger")

	err = handler.BookmarkDeleteHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
}

func TestBookmarkEditViewHandler(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	bookmark, err := repo.Create(data.BookmarkForm{
		URL:   "https://example.com/public",
		Title: "Example public",
	})
	require.NoError(t, err)

	e := router.NewBaseApp(db)

	// view authenticated
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/bookmarks/%d/edit", bookmark.ID), strings.NewReader(""))
	rec := httptest.NewRecorder()
	sc := test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues(fmt.Sprint(bookmark.ID))

	err = handler.BookmarkEditViewHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Contains(t, rec.Body.String(), bookmark.Title)

	// view unauthenticated
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/bookmarks/%d/edit", bookmark.ID), strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewUnauthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues(fmt.Sprint(bookmark.ID))

	err = handler.BookmarkEditViewHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.True(t, strings.HasPrefix(rec.Result().Header.Get("Location"), "/login?next="))

	// view non-existing bookmark
	req = httptest.NewRequest(http.MethodGet, "/bookmarks/42/edit", strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues("42")

	err = handler.BookmarkEditViewHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)

	// view non-existing bookmark with non-integer id
	req = httptest.NewRequest(http.MethodGet, "/bookmarks/notaninteger/edit", strings.NewReader(""))
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues("notaninteger")

	err = handler.BookmarkEditViewHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
}

func TestBookmarkEditHandler(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	bookmark, err := repo.Create(data.BookmarkForm{
		URL:   "https://example.com/private",
		Title: "Example public",
	})
	require.NoError(t, err)

	e := router.NewBaseApp(db)

	contentType := "application/x-www-form-urlencoded"

	// success
	form := url.Values{}
	form.Add("url", "https://example.com/public")
	form.Add("title", "Example - About")
	form.Add("description", "About example.com")
	form.Add("public", "on")

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/bookmarks/%d/edit", bookmark.ID), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()
	sc := test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues(fmt.Sprint(bookmark.ID))

	err = handler.BookmarkEditHandler(sc)
	require.NoError(t, err)

	var actual data.Bookmark
	result := db.Last(&actual)
	require.NoError(t, result.Error)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.Equal(t, fmt.Sprintf("/bookmarks/%d", actual.ID), rec.Header().Get("Location"))

	require.Equal(t, form.Get("url"), actual.URL)
	require.Equal(t, form.Get("title"), actual.Title)
	require.Equal(t, form.Get("description"), actual.Description)
	require.Equal(t, data.BookmarkPrivacyPublic, actual.Privacy)

	// validation error
	form = url.Values{}
	form.Add("url", "")
	form.Add("title", "Example - About")

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/bookmarks/%d/edit", bookmark.ID), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues(fmt.Sprint(bookmark.ID))

	err = handler.BookmarkEditHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Contains(t, rec.Body.String(), "Failed to edit bookmark")
	require.Contains(t, rec.Body.String(), "URL is required")

	// edit non-existing bookmark
	form = url.Values{}
	form.Add("url", "https://example.com/public")
	form.Add("title", "Example - About")
	form.Add("description", "About example.com")
	form.Add("public", "on")

	req = httptest.NewRequest(http.MethodPost, "/bookmarks/42/edit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues("42")

	err = handler.BookmarkEditHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)

	// edit non-existing bookmark with non-integer id
	form = url.Values{}
	form.Add("url", "https://example.com/public")
	form.Add("title", "Example - About")
	form.Add("description", "About example.com")
	form.Add("public", "on")

	req = httptest.NewRequest(http.MethodPost, "/bookmarks/notaninteger/edit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)
	rec = httptest.NewRecorder()
	sc = test.NewAuthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues("notaninteger")

	err = handler.BookmarkEditHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)

	// unauthenticated
	form = url.Values{}
	form.Add("url", "https://example.com/public")
	form.Add("title", "Example - About")
	form.Add("description", "About example.com")
	form.Add("public", "on")

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/bookmarks/%d/edit", bookmark.ID), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)
	rec = httptest.NewRecorder()
	sc = test.NewUnauthenticatedContext(e.NewContext(req, rec), db)
	sc.SetParamNames("id")
	sc.SetParamValues(fmt.Sprint(bookmark.ID))

	err = handler.BookmarkEditHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.Equal(t, "/login", rec.Header().Get("Location"))
}
