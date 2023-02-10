package handler

import (
	_ "embed"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/middleware"
	"github.com/chdorner/submarine/util"
	"github.com/labstack/echo/v4"
)

func LoginViewHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)

	next, _ := url.QueryUnescape(sc.QueryParam("next"))

	return sc.Render(http.StatusOK, "login.html", map[string]interface{}{
		"next": next,
	})
}

func LoginHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)

	next := sc.FormValue("next")
	genericError := map[string]interface{}{
		"error": "Login failed!",
		"next":  next,
	}

	settingsRepo := data.NewSettingsRepository(sc.DB)
	settings, err := settingsRepo.Get()
	if err != nil || settings == nil {
		return sc.Render(http.StatusOK, "login.html", genericError)
	}

	if !util.ComparePassword(sc.FormValue("password"), settings.Password) {
		return sc.Render(http.StatusOK, "login.html", genericError)
	}

	repo := data.NewSessionRepository(sc.DB)
	session, err := repo.Create(&data.SessionCreate{
		IP:        sc.RealIP(),
		UserAgent: sc.Request().UserAgent(),
	})
	if err != nil {
		return sc.Render(http.StatusOK, "login.html", map[string]interface{}{
			"error": "Failed to create session!",
			"next":  next,
		})
	}
	middleware.SetCookieSessionToken(sc, session.Token)

	if next != "" {
		next, err := base64.StdEncoding.DecodeString(next)
		if err == nil {
			return sc.Redirect(http.StatusFound, string(next))
		}
	}

	return sc.Redirect(http.StatusFound, "/")
}

func LogoutHandler(c echo.Context) error {
	sc := c.(*middleware.SubmarineContext)

	middleware.ClearCookieSessionToken(c)

	return sc.Redirect(http.StatusFound, "/")
}
