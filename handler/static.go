package handler

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed static
	staticFS  embed.FS
	checksums map[string]string
)

func NewStaticHandler() (echo.HandlerFunc, error) {
	err := calculateChecksums()
	if err != nil {
		return nil, err
	}

	fileHandler := http.FileServer(http.FS(staticFS))
	return func(c echo.Context) error {
		response := c.Response()

		if c.QueryParam("v") != "" {
			response.Header().Add("Cache-Control", "max-age=315360000")
		}

		fileHandler.ServeHTTP(response, c.Request())
		return nil
	}, nil
}

func StaticAssetPath(name string) string {
	path := fmt.Sprintf("/static/%s", name)
	checksum, ok := checksums[name]
	if ok {
		path = fmt.Sprintf("%s?v=%s", path, checksum)
	}
	return path
}

func calculateChecksums() error {
	checksums = make(map[string]string)
	entries, err := staticFS.ReadDir("static")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		contents, err := staticFS.ReadFile(path.Join("static", entry.Name()))
		if err != nil {
			return err
		}

		h := sha256.New()
		_, err = io.Copy(h, bytes.NewReader(contents))
		if err != nil {
			return err
		}
		checksums[entry.Name()] = fmt.Sprintf("%x", h.Sum(nil))
	}

	return nil
}
