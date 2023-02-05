package web

import (
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from submarine")) //nolint:errcheck
}
