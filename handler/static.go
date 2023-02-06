package handler

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed static
	staticFS embed.FS
)

func NewStaticHandler() (echo.HandlerFunc, error) {
	handler := http.FileServer(http.FS(staticFS))
	return echo.WrapHandler(handler), nil
}
