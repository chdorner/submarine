package data_test

import (
	"fmt"
	"testing"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/test"
	"github.com/stretchr/testify/require"
)

func TestBookmarkRepositoryGet(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	expected := data.Bookmark{URL: "https://example.com"}
	result := db.Create(&expected)
	require.NoError(t, result.Error)

	actual, err := repo.Get(expected.ID)
	require.NoError(t, err)
	require.Equal(t, expected.ID, actual.ID)

	// not found
	actual, err = repo.Get(0)
	require.NoError(t, err)
	require.Nil(t, actual)
}

func TestBookmarkRepositoryCreate(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	req := data.BookmarkForm{
		URL:         "https://en.wikipedia.org/wiki/Main_Page",
		Title:       "Wikipedia",
		Description: "the free encyclopedia that anyone can edit",
		Public:      true,
		Tags:        "to-read, articles",
	}
	bookmark, err := repo.Create(req)
	require.NoError(t, err)
	require.Equal(t, req.URL, bookmark.URL)
	require.Equal(t, req.Title, bookmark.Title)
	require.Equal(t, req.Description, bookmark.Description)
	require.Equal(t, data.BookmarkPrivacyPublic, bookmark.Privacy)
	tagNames := []string{}
	for _, tag := range bookmark.Tags {
		tagNames = append(tagNames, tag.DisplayName)
	}
	require.Equal(t, []string{"to-read", "articles"}, tagNames)

	// minimum required fields
	_, err = repo.Create(data.BookmarkForm{
		URL: "https://en.wikipedia.org/wiki/Main_Page",
	})
	require.NoError(t, err)

	// requires URL
	_, err = repo.Create(data.BookmarkForm{})
	require.ErrorContains(t, err, "bookmarks.url")
}

func TestBookmarkRepositoryList(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	for i := 0; i < 25; i++ {
		public := i%2 == 0
		tag := "private"
		if public {
			tag = "public"
		}
		_, err := repo.Create(data.BookmarkForm{
			URL:    fmt.Sprintf("https://example-%d.com", i),
			Title:  fmt.Sprintf("Bookmark %d", i),
			Public: public,
			Tags:   tag,
		})
		require.NoError(t, err)
	}

	// privacy all - page 1
	result, err := repo.List(data.BookmarkListRequest{
		Privacy:              data.BookmarkPrivacyQueryAll,
		PaginationPathPrefix: "/?",
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 10)
	require.Equal(t, "Bookmark 0", result.Items[0].Title)
	require.Equal(t, "Bookmark 9", result.Items[len(result.Items)-1].Title)
	require.Equal(t, int64(25), result.Count)
	require.False(t, result.HasPrev)
	require.True(t, result.HasNext)
	require.Equal(t, "/?offset=10", result.NextURL)

	// privacy all - page 2
	result, err = repo.List(data.BookmarkListRequest{
		Privacy: data.BookmarkPrivacyQueryAll,
		Offset:  10,

		PaginationPathPrefix: "/?",
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 10)
	require.Equal(t, "Bookmark 10", result.Items[0].Title)
	require.Equal(t, "Bookmark 19", result.Items[len(result.Items)-1].Title)
	require.Equal(t, int64(25), result.Count)
	require.True(t, result.HasPrev)
	require.Equal(t, "/?offset=0", result.PrevURL)
	require.True(t, result.HasNext)
	require.Equal(t, "/?offset=20", result.NextURL)

	// privacy all - page 3
	result, err = repo.List(data.BookmarkListRequest{
		Privacy:              data.BookmarkPrivacyQueryAll,
		Offset:               20,
		PaginationPathPrefix: "/?",
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 5)
	require.Equal(t, "Bookmark 20", result.Items[0].Title)
	require.Equal(t, "Bookmark 24", result.Items[len(result.Items)-1].Title)
	require.Equal(t, int64(25), result.Count)
	require.True(t, result.HasPrev)
	require.Equal(t, "/?offset=10", result.PrevURL)
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

	// privacy all, filter tag
	tagRepo := data.NewTagRepository(db)
	tag, err := tagRepo.GetByName("public")
	require.NoError(t, err)
	result, err = repo.List(data.BookmarkListRequest{
		Privacy: data.BookmarkPrivacyPublic,
		TagID:   tag.ID,
	})
	require.NoError(t, err)
	require.Equal(t, int64(13), result.Count)
	for _, item := range result.Items {
		require.Len(t, item.Tags, 1)
		require.Equal(t, tag.ID, item.Tags[0].ID)
	}
}

func TestBookmarkRepositorySearch(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	for i := 0; i < 25; i++ {
		_, err := repo.Create(data.BookmarkForm{
			URL:         fmt.Sprintf("https://example-%d.com", i),
			Title:       fmt.Sprintf("Bookmark %d", i),
			Description: fmt.Sprintf("Bookmark Description %d", i),
		})
		require.NoError(t, err)
	}
	_, err := repo.Create((data.BookmarkForm{
		URL:         "https://other.org",
		Title:       "Other",
		Description: "Other Description",
	}))
	require.NoError(t, err)

	// search bookmark
	result, err := repo.Search(data.BookmarkSearchRequest{Query: "bookmark"})
	require.NoError(t, err)
	require.Len(t, result.Items, 10)
	require.Equal(t, int64(25), result.Count)
	require.True(t, result.HasNext)
	require.Equal(t, "offset=10", result.NextURL)
	require.False(t, result.HasPrev)

	// search description
	result, err = repo.Search(data.BookmarkSearchRequest{Query: "description"})
	require.NoError(t, err)
	require.Equal(t, int64(26), result.Count)

	// search second page
	result, err = repo.Search(data.BookmarkSearchRequest{
		Query:  "bookmark",
		Offset: 10,
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 10)
	require.True(t, result.HasNext)
	require.Equal(t, "offset=20", result.NextURL)
	require.True(t, result.HasPrev)
	require.Equal(t, "offset=0", result.PrevURL)

	// search third page
	result, err = repo.Search(data.BookmarkSearchRequest{
		Query:  "bookmark",
		Offset: 20,
	})
	require.NoError(t, err)
	require.Len(t, result.Items, 5)
	require.False(t, result.HasNext)
	require.True(t, result.HasPrev)
	require.Equal(t, "offset=10", result.PrevURL)

	// search no result
	result, err = repo.Search(data.BookmarkSearchRequest{Query: "nothing"})
	require.NoError(t, err)
	require.Len(t, result.Items, 0)
	require.Equal(t, int64(0), result.Count)
}

func TestBookmarkRepositoryDelete(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	expected := data.Bookmark{URL: "https://example.com"}
	result := db.Create(&expected)
	require.NoError(t, result.Error)

	err := repo.Delete(expected.ID)
	require.NoError(t, err)
	deleted, err := repo.Get(expected.ID)
	require.NoError(t, err)
	require.Nil(t, deleted)

	// not found
	err = repo.Delete(42)
	require.EqualError(t, err, "bookmark with id 42 not found")
}

func TestBookmarkRepositoryUpdate(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewBookmarkRepository(db)

	bookmark, err := repo.Create(data.BookmarkForm{
		URL:         "https://en.wikipedia.org",
		Title:       "English Wikipedia",
		Description: "English language",
		Public:      true,
		Tags:        "toRead, articles",
	})
	require.NoError(t, err)
	require.NotNil(t, bookmark)

	expected := data.BookmarkForm{
		URL:         "https://de.wikipedia.org",
		Title:       "German Wikipedia",
		Description: "German language",
		Public:      false,
		Tags:        "articles, recommended, top10",
	}
	err = repo.Update(bookmark.ID, expected)
	require.NoError(t, err)

	actual, err := repo.Get(bookmark.ID)
	require.NoError(t, err)
	require.Equal(t, expected.URL, actual.URL)
	require.Equal(t, expected.Title, actual.Title)
	require.Equal(t, expected.Description, actual.Description)
	require.False(t, actual.IsPublic())
	require.Len(t, actual.Tags, 3)
	require.Equal(t, "articles", actual.Tags[0].DisplayName)
	require.Equal(t, "recommended", actual.Tags[1].DisplayName)
	require.Equal(t, "top10", actual.Tags[2].DisplayName)

	// update not found
	err = repo.Update(uint(42), expected)
	require.EqualError(t, err, "bookmark with id 42 not found")
}
