package test

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/chdorner/submarine/middleware"
)

func NewAuthenticatedContext(c echo.Context, db *gorm.DB) *middleware.SubmarineContext {
	sc := middleware.InitSubmarineContext(c, db)
	sc.Set("IsAuthenticated", true)
	return sc
}

func NewUnauthenticatedContext(c echo.Context, db *gorm.DB) *middleware.SubmarineContext {
	return middleware.InitSubmarineContext(c, db)
}
