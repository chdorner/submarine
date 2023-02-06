package handler_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/chdorner/submarine/handler"
	"github.com/chdorner/submarine/handler/testdata"
)

func TestTemplatesParse(t *testing.T) {
	tpl := &handler.Templates{}

	err := tpl.Parse(testdata.ViewTemplateFiles, testdata.CommonTemplateFiles)
	require.NoError(t, err)

	rendered := bytes.NewBuffer([]byte{})
	err = tpl.Registry["view.html"].ExecuteTemplate(rendered, "base", nil)
	require.NoError(t, err)

	require.Equal(t, "\ntest-base\n\ntest-view\n\n", rendered.String())
}
