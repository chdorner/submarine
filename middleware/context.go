package middleware

import "github.com/labstack/echo/v4"

type SubmarineContext struct {
	echo.Context
	IsAuthenticated bool
	SessionID       string
}

func SubmarineContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(InitSubmarineContext(c))
	}
}

func InitSubmarineContext(c echo.Context) *SubmarineContext {
	return &SubmarineContext{
		c,
		false, // IsAuthenticated
		"",    // SessionID
	}
}
