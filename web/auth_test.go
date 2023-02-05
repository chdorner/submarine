package web_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chdorner/submarine/web"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestCookieAuthMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(nil))
	req.AddCookie(&http.Cookie{
		Name:  "SubmarineSessionID",
		Value: "the-session-id",
	})
	rec := httptest.NewRecorder()
	sc := web.InitSubmarineContext(e.NewContext(req, rec))

	var actual string
	handler := func(c echo.Context) error {
		sc := c.(*web.SubmarineContext)
		actual = sc.SessionID
		return c.String(http.StatusOK, "OK")
	}

	err := web.CookieAuthMiddleware(handler)(sc)
	require.NoError(t, err)
	require.Equal(t, "the-session-id", actual)
}

func TestSetCookieSessionID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(nil))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error {
		web.SetCookieSessionID(c, "the-session-id")
		return c.String(http.StatusOK, "Successfully logged-in")
	}

	err := handler(c)
	require.NoError(t, err)
	cookie := rec.Header().Get("Set-Cookie")
	require.Contains(t, cookie, "SubmarineSessionID=the-session-id")
}
