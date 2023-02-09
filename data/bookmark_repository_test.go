package data_test

import (
	"fmt"
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
	require.Equal(t, req.Description, bookmark.Description)

	// minimum required fields
	_, err = repo.Create(data.BookmarkCreate{
		URL: "https://en.wikipedia.org/wiki/Main_Page",
	})
	require.NoError(t, err)

	// requires URL
	_, err = repo.Create(data.BookmarkCreate{})
	require.ErrorContains(t, err, "bookmarks.url")
}

func TestBookmarkRepositoryList(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	for i := 0; i < 25; i++ {
		privacy := data.BookmarkPrivacyPrivate
		if i%2 == 0 {
			privacy = data.BookmarkPrivacyPublic
		}
		_, err := repo.Create(data.BookmarkCreate{
			URL:     fmt.Sprintf("https://example-%d.com", i),
			Title:   fmt.Sprintf("Bookmark %d", i),
			Privacy: privacy,
		})
		require.NoError(t, err)
	}

	// privacy all - page 1
	result, err := repo.List(data.BookmarkListRequest{
		Privacy: data.BookmarkPrivacyQueryAll,
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 10)
	require.Equal(t, "Bookmark 0", result.Items[0].Title)
	require.Equal(t, "Bookmark 9", result.Items[len(result.Items)-1].Title)
	require.Equal(t, int64(25), result.Count)
	require.False(t, result.HasPrev)
	require.True(t, result.HasNext)
	require.Equal(t, 10, result.NextOffset)

	// privacy all - page 2
	result, err = repo.List(data.BookmarkListRequest{
		Privacy: data.BookmarkPrivacyQueryAll,
		Offset:  10,
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 10)
	require.Equal(t, "Bookmark 10", result.Items[0].Title)
	require.Equal(t, "Bookmark 19", result.Items[len(result.Items)-1].Title)
	require.Equal(t, int64(25), result.Count)
	require.True(t, result.HasPrev)
	require.Equal(t, 0, result.PrevOffset)
	require.True(t, result.HasNext)
	require.Equal(t, 20, result.NextOffset)

	// privacy all - page 3
	result, err = repo.List(data.BookmarkListRequest{
		Privacy: data.BookmarkPrivacyQueryAll,
		Offset:  20,
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 5)
	require.Equal(t, "Bookmark 20", result.Items[0].Title)
	require.Equal(t, "Bookmark 24", result.Items[len(result.Items)-1].Title)
	require.Equal(t, int64(25), result.Count)
	require.True(t, result.HasPrev)
	require.Equal(t, 10, result.PrevOffset)
	require.False(t, result.HasNext)

	// privacy all - order
	result, err = repo.List(data.BookmarkListRequest{
		Privacy: data.BookmarkPrivacyQueryAll,
		Order:   "created_at desc",
	})
	require.NoError(t, err)
	require.Equal(t, "Bookmark 24", result.Items[0].Title)

	// privacy private
	result, err = repo.List(data.BookmarkListRequest{
		Privacy: data.BookmarkPrivacyPrivate,
	})
	require.NoError(t, err)
	require.Equal(t, int64(12), result.Count)
	for _, item := range result.Items {
		require.Equal(t, data.BookmarkPrivacyPrivate, item.Privacy)
	}

	// privacy public
	result, err = repo.List(data.BookmarkListRequest{
		Privacy: data.BookmarkPrivacyPublic,
	})
	require.NoError(t, err)
	require.Equal(t, int64(13), result.Count)
	for _, item := range result.Items {
		require.Equal(t, data.BookmarkPrivacyPublic, item.Privacy)
	}

	// privacy defaulting to public
	result, err = repo.List(data.BookmarkListRequest{})
	require.NoError(t, err)
	require.Equal(t, int64(13), result.Count)
	for _, item := range result.Items {
		require.Equal(t, data.BookmarkPrivacyPublic, item.Privacy)
	}
}
