package data_test

import (
	"testing"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/test"
	"github.com/stretchr/testify/require"
)

func TestBookmarkRepositoryCreate(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	req := data.BookmarkCreate{
		URL:         "https://en.wikipedia.org/wiki/Main_Page",
		Title:       "Wikipedia",
		Description: "the free encyclopedia that anyone can edit",
	}
	bookmark, err := repo.Create(req)
	require.NoError(t, err)
	require.Equal(t, req.URL, bookmark.URL)
	require.Equal(t, req.Title, bookmark.Title)
	require.Equal(t, req.Description, bookmark.Descripton)

	// minimum required fields
	_, err = repo.Create(data.BookmarkCreate{
		URL: "https://en.wikipedia.org/wiki/Main_Page",
	})
	require.NoError(t, err)

	// requires URL
	_, err = repo.Create(data.BookmarkCreate{})
	require.ErrorContains(t, err, "bookmarks.url")
}
