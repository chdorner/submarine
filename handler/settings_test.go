package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chdorner/submarine/handler"
	"github.com/chdorner/submarine/router"
	"github.com/chdorner/submarine/test"
	"github.com/stretchr/testify/require"
)

func TestSettingsHandler(t *testing.T) {
	// unauthenticated
	e := router.NewBaseApp(nil)
	req := httptest.NewRequest(http.MethodGet, "/settings", strings.NewReader(""))
	rec := httptest.NewRecorder()
	sc := test.NewUnauthenticatedContext(e.NewContext(req, rec), nil)
	err := handler.SettingsHandler(sc)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Result().StatusCode)
	require.True(t, strings.HasPrefix(rec.Result().Header.Get("Location"), "/login?next="))
}
