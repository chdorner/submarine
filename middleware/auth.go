package middleware

import (
	"net/http"
	"time"

	"github.com/chdorner/submarine/data"
	"github.com/labstack/echo/v4"
)

func CookieAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sc := c.(*SubmarineContext)

		token := getCookieSessionToken(sc)
		if token != "" {
			repo := data.NewSessionRepository(sc.DB)
			session, err := repo.GetByToken(token)
			if err == nil {
				sc.Set("SessionID", session.ID)
				sc.Set("IsAuthenticated", true)
			}
		}

		err := next(sc)
		if err != nil {
			c.Error(err)
		}

		return nil
	}
}

func getCookieSessionToken(c echo.Context) string {
	cookie, err := c.Cookie("SubmarineSessionToken")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func SetCookieSessionToken(c echo.Context, token string) {
	setCookie(c, &http.Cookie{
		Name:    "SubmarineSessionToken",
		Value:   token,
		Expires: time.Now().Add(14 * 24 * time.Hour),
	})
}

func ClearCookieSessionToken(c echo.Context) {
	setCookie(c, &http.Cookie{
		Name:    "SubmarineSessionToken",
		Value:   "",
		Expires: time.Now().Add(-time.Second),
	})
}

func setCookie(c echo.Context, cookie *http.Cookie) {
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Secure = c.Scheme() == "https"
	c.SetCookie(cookie)
}
