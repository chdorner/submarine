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

func NewBaseApp(db *gorm.DB) *echo.Echo {
	e := echo.New()
	e.Pre(echomiddleware.RemoveTrailingSlashWithConfig(echomiddleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Renderer = handler.NewTemplates()

	e.Use(echomiddleware.Recover())
	e.Use(middleware.SubmarineContextMiddleware(db))
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
	e.Use(echomiddleware.CSRFWithConfig(echomiddleware.CSRFConfig{
		TokenLookup: "form:_csrf",
	}))
	e.Use(middleware.CookieAuthMiddleware)

	return e
}

func New(db *gorm.DB) *echo.Echo {
	e := NewBaseApp(db)

	e.GET("/", handler.BookmarksListHandler)
	e.POST("/bookmarks/:id/delete", handler.BookmarkDeleteHandler)
	e.GET("/bookmarks/:id/edit", handler.BookmarkEditViewHandler)
	e.POST("/bookmarks/:id/edit", handler.BookmarkEditHandler)
	e.GET("/bookmarks/:id", handler.BookmarkShowHandler)
	e.GET("/bookmarks/new", handler.BookmarksNewHandler)
	e.POST("/bookmarks", handler.BookmarksCreateHandler)

	e.GET("/tags/:name", handler.TagHandler)

	e.GET("/settings", handler.SettingsHandler)

	e.GET("/login", handler.LoginViewHandler)
	e.POST("/login", handler.LoginHandler)
	e.GET("/logout", handler.LogoutHandler)

	staticHandler, err := handler.NewStaticHandler()
	if err != nil {
		panic(err)
	}
	e.GET("/static/*", staticHandler)

	return e
}
