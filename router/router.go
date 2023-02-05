package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/chdorner/submarine/handler"
	"github.com/chdorner/submarine/middleware"
)

func New(db *gorm.DB) *echo.Echo {
	e := echo.New()
	e.Pre(echomiddleware.RemoveTrailingSlashWithConfig(echomiddleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Renderer = handler.NewTemplates()

	e.GET("/", handler.RootHandler)
	e.GET("/login", handler.LoginViewHandler)

	// Middleware
	e.Use(middleware.SubmarineContextMiddleware)
	log := logrus.New()
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.WithFields(logrus.Fields{
					"URI":    v.URI,
					"status": v.Status,
				}).Info("request")
			} else {
				log.WithFields(logrus.Fields{
					"URI":    v.URI,
					"status": v.Status,
					"error":  v.Error,
				}).Error("request error")
			}
			return nil
		},
	}))
	e.Use(echomiddleware.Recover())
	e.Use(middleware.CookieAuthMiddleware(db))

	return e
}
