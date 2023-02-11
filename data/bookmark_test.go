package data_test

import (
	"fmt"
	"testing"

	"github.com/chdorner/submarine/data"
	"github.com/stretchr/testify/require"
)

func TestBookmarkCreateIsValid(t *testing.T) {
	req := data.BookmarkCreate{URL: "https://en.wikipedia.org/wiki/Main_Page"}
	err := req.IsValid()
	require.Nil(t, err)

	// empty URL
	req = data.BookmarkCreate{URL: ""}
	err = req.IsValid()
	require.Error(t, err)
	require.Equal(t, "URL is required", err.Fields["URL"])

	// invalid URL
	invalidURLCases := []string{
		"/path",
		"http:",
		"https://",
		"host/path",
	}

	for _, tc := range invalidURLCases {
		t.Run(fmt.Sprintf("invalid URL %s", tc), func(t *testing.T) {
			req = data.BookmarkCreate{URL: tc}
			err = req.IsValid()
			require.Error(t, err)
			require.Equal(t, "URL format is invalid", err.Fields["URL"])
		})
	}

	// invalid Privacy
	req = data.BookmarkCreate{
		URL:     "https://example.com",
		Privacy: data.BookmarkPrivacyQueryAll,
	}
	err = req.IsValid()
	require.Error(t, err)
	require.Equal(t, "Invalid permission", err.Fields["Privacy"])
}
