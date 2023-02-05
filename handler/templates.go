package handler

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed templates
	embeddedFiles embed.FS
)

type Templates struct {
	templates *template.Template
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseFS(embeddedFiles, "templates/*.html")),
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
