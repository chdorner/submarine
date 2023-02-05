package web

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func CookieAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sc := c.(*SubmarineContext)
		sc.SessionID = getCookieSessionID(sc)

		err := next(sc)
		if err != nil {
			c.Error(err)
		}

		return nil
	}
}

func getCookieSessionID(c echo.Context) string {
	cookie, err := c.Cookie("SubmarineSessionID")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func SetCookieSessionID(c echo.Context, sessionID string) {
	cookie := &http.Cookie{
		Name:     "SubmarineSessionID",
		Value:    sessionID,
		SameSite: http.SameSiteStrictMode,
		// TODO: enable secure cookie when using HTTPS
		//Secure: true,
		Expires: time.Now().Add(14 * 24 * time.Hour),
	}
	c.SetCookie(cookie)
}
