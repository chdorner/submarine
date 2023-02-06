package testdata

import "embed"

var (
	//go:embed templates/common/*.html
	CommonTemplateFiles embed.FS

	//go:embed templates/views/*.html
	ViewTemplateFiles embed.FS
)
