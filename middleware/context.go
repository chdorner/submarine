package middleware

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SubmarineContext struct {
	echo.Context
	DB              *gorm.DB
	IsAuthenticated bool
	SessionID       uint
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
		db,    // DB
		false, // IsAuthenticated
		0,     // SessionID
	}
}
