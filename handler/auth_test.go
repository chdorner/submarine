package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/handler"
	"github.com/chdorner/submarine/middleware"
	"github.com/chdorner/submarine/router"
	"github.com/chdorner/submarine/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestLoginHandler(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()

	contentType := "application/x-www-form-urlencoded"

	e := router.NewBaseApp(db)
	form := strings.NewReader("password=secret")
	req := httptest.NewRequest(http.MethodPost, "/login", form)
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()
	sc := middleware.InitSubmarineContext(e.NewContext(req, rec), db)

	// no site settings available
	err := handler.LoginHandler(sc)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), "Login failed!")

	repo := data.NewSettingsRepository(db)
	err = repo.Upsert(data.SettingsUpsert{
		Password: "secret",
	})
	require.NoError(t, err)

	// with incorrect password
	form = strings.NewReader("password=supersecret")
	req = httptest.NewRequest(http.MethodPost, "/login", form)
	req.Header.Set("Content-Type", contentType)
	rec = httptest.NewRecorder()
	sc = middleware.InitSubmarineContext(e.NewContext(req, rec), db)
	err = handler.LoginHandler(sc)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), "Login failed!")

	// with correct password
	form = strings.NewReader("password=secret")
	req = httptest.NewRequest(http.MethodPost, "/login", form)
	req.Header.Set("Content-Type", contentType)
	rec = httptest.NewRecorder()
	sc = middleware.InitSubmarineContext(e.NewContext(req, rec), db)
	err = handler.LoginHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.Equal(t, "/", rec.Header().Get("Location"))
	cookie := test.ParseCookie(t, rec.Header().Get("Set-Cookie"))
	uuid.MustParse(cookie["SubmarineSessionToken"].(string))
}

func TestLogoutHandler(t *testing.T) {
	e := router.NewBaseApp(nil)
	req := httptest.NewRequest(http.MethodGet, "/logout", strings.NewReader(""))
	rec := httptest.NewRecorder()
	sc := middleware.InitSubmarineContext(e.NewContext(req, rec), nil)
	err := handler.LogoutHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.Equal(t, "/", rec.Header().Get("Location"))
	cookie := test.ParseCookie(t, rec.Header().Get("Set-Cookie"))
	require.Empty(t, cookie["SubmarineSessionToken"])
}
