package test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func ParseCookie(t *testing.T, value string) map[string]interface{} {
	cookie := map[string]interface{}{}
	for _, token := range strings.Split(value, ";") {
		tokenSplits := strings.Split(token, "=")
		cookie[strings.TrimSpace(tokenSplits[0])] = tokenSplits[1]
	}

	expires, ok := cookie["Expires"]
	if ok {
		parsed, err := time.Parse(time.RFC1123, expires.(string))
		require.NoError(t, err)
		cookie["Expires"] = parsed
	}

	return cookie
}
