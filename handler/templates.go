package handler

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

var (
	//go:embed templates/common/*.html
	commonTemplateFiles embed.FS

	//go:embed templates/views/*.html
	viewTemplateFiles embed.FS
)

type Templates struct {
	Registry map[string]*template.Template
}

func NewTemplates() *Templates {
	t := &Templates{}
	err := t.Parse(viewTemplateFiles, commonTemplateFiles)
	if err != nil {
		panic(err)
	}
	return t
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tpl, ok := t.Registry[name]
	if !ok {
		panic(fmt.Sprintf("template %s does not exists", name))
	}

	funcMap := template.FuncMap{
		"IsAuthenticated": func() bool {
			IsAuthenticated := c.Get("IsAuthenticated")
			if IsAuthenticated == nil {
				return false
			}
			return IsAuthenticated.(bool)
		},
		"CSRFHiddenInput": func() template.HTML {
			token := c.Get(echomiddleware.DefaultCSRFConfig.ContextKey)
			if token == nil {
				return ""
			}
			return template.HTML(
				fmt.Sprintf(`<input type="hidden" name="_csrf" value="%s">`, token),
			)
		},
	}

	dataMap := map[string]interface{}{}
	if data != nil {
		dataMap = data.(map[string]interface{})
	}

	dataMap["IsAuthenticated"] = c.Get("IsAuthenticated")

	return tpl.Funcs(funcMap).ExecuteTemplate(w, "base", dataMap)
}

func (t *Templates) Parse(views, common fs.ReadDirFS) error {
	t.Registry = make(map[string]*template.Template)

	// empty functions need to be overwritten when executing templates
	funcMap := template.FuncMap{
		"StaticAssetPath": func(name string) template.HTML {
			return template.HTML(StaticAssetPath(name))
		},
		"IsAuthenticated": func() bool {
			return false
		},
		"CSRFHiddenInput": func() template.HTML {
			return template.HTML("")
		},
	}

	var commonContents strings.Builder
	entries, err := common.ReadDir(path.Join("templates", "common"))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		data, err := fs.ReadFile(common, path.Join("templates", "common", entry.Name()))
		if err != nil {
			return err
		}
		commonContents.Write(data)
	}

	entries, err = views.ReadDir(path.Join("templates", "views"))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		data, err := fs.ReadFile(views, path.Join("templates", "views", entry.Name()))
		if err != nil {
			return err
		}

		var templateContents strings.Builder
		templateContents.WriteString(commonContents.String())
		templateContents.Write(data)

		t.Registry[entry.Name()] = template.Must(template.New("main").Funcs(funcMap).Parse(templateContents.String()))
	}

	return nil
}
