package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/middleware"
	"github.com/chdorner/submarine/test"
)

func TestCookieAuthMiddleware(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()

	repo := data.NewSessionRepository(db)
	session, err := repo.Create(&data.SessionCreate{})
	require.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(nil))
	req.AddCookie(&http.Cookie{
		Name:  "SubmarineSessionToken",
		Value: session.Token,
	})
	rec := httptest.NewRecorder()
	sc := middleware.InitSubmarineContext(e.NewContext(req, rec), db)

	var actualID uint
	var actualAuthenticated bool
	handler := func(c echo.Context) error {
		sc := c.(*middleware.SubmarineContext)
		actualID = sc.Get("SessionID").(uint)
		actualAuthenticated = sc.IsAuthenticated()
		return c.String(http.StatusOK, "OK")
	}

	err = middleware.CookieAuthMiddleware(handler)(sc)
	require.NoError(t, err)
	require.True(t, actualAuthenticated)
	require.Equal(t, session.ID, actualID)
}

func TestSetCookieSessionToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(nil))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error {
		middleware.SetCookieSessionToken(c, "the-session-token")
		return c.String(http.StatusOK, "Successfully logged-in")
	}

	err := handler(c)
	require.NoError(t, err)
	cookie := test.ParseCookie(t, rec.Header().Get("Set-Cookie"))
	require.Equal(t, cookie["SubmarineSessionToken"], "the-session-token")
	require.Greater(t, cookie["Expires"], time.Now())
}

func TestClearCookieSessionToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(nil))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error {
		middleware.ClearCookieSessionToken(c)
		return c.String(http.StatusOK, "Successfully logged-out")
	}

	err := handler(c)
	require.NoError(t, err)
	cookie := test.ParseCookie(t, rec.Header().Get("Set-Cookie"))
	require.Empty(t, cookie["SubmarineSessionToken"])
	require.Less(t, cookie["Expires"], time.Now())
}
