package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SubmarineContext struct {
	echo.Context
	DB *gorm.DB
}

func SubmarineContextMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(InitSubmarineContext(c, db))
		}
	}
}

func InitSubmarineContext(c echo.Context, db *gorm.DB) *SubmarineContext {
	return &SubmarineContext{
		c,
		db, // DB
	}
}

func (sc *SubmarineContext) IsAuthenticated() bool {
	isAuthenticated := sc.Get("IsAuthenticated")
	if isAuthenticated == nil {
		return false
	}
	return isAuthenticated.(bool)
}

func (sc *SubmarineContext) RedirectToLogin() error {
	redirect := "/login"
	if sc.Request().Method == http.MethodGet {
		next := sc.Path()
		query := sc.QueryString()
		if query != "" {
			next = fmt.Sprintf("%s?%s", next, query)
		}
		nextEncoded := base64.StdEncoding.EncodeToString([]byte(next))
		redirect = fmt.Sprintf("%s?next=%s", redirect, nextEncoded)
	}
	return sc.Redirect(http.StatusFound, redirect)
}
