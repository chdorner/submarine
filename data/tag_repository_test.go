package data_test

import (
	"testing"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/test"
	"github.com/stretchr/testify/require"
)

func TestTagRepositoryGet(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewTagRepository(db)

	expected := data.Tag{Name: "toRead"}
	result := db.Create(&expected)
	require.NoError(t, result.Error)

	actual, err := repo.GetByName(expected.Name)
	require.NoError(t, err)
	require.Equal(t, expected.ID, actual.ID)

	// finds case insensitive
	actual, err = repo.GetByName("TOREAD")
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected.ID, actual.ID)

	// not found
	actual, err = repo.GetByName("missing")
	require.NoError(t, err)
	require.Nil(t, actual)

	// empty name
	actual, err = repo.GetByName("")
	require.NoError(t, err)
	require.Nil(t, actual)
}

func TestTagRepositoryUpsert(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewTagRepository(db)

	// create toRead and articles
	tagNames := []string{"toRead", "articles"}
	var count int64
	db.Model(&data.Tag{}).Where("name IN ?", tagNames).Count(&count)
	require.Equal(t, int64(0), count)
	tags, err := repo.Upsert(tagNames)
	require.NoError(t, err)
	require.Equal(t, "toRead", tags[0].Name)
	require.Equal(t, "articles", tags[1].Name)
	db.Model(&data.Tag{}).Count(&count)
	require.Equal(t, int64(2), count)

	// update articles and create recommended
	db.Model(&data.Tag{}).Where("name = ?", "recommended").Count(&count)
	require.Equal(t, int64(0), count)
	tags, err = repo.Upsert([]string{"articles", "recommended"})
	require.NoError(t, err)
	db.Model(&data.Tag{}).Count(&count)
	require.Equal(t, int64(3), count)
	require.Equal(t, "articles", tags[0].Name)
	require.Equal(t, "recommended", tags[1].Name)

	// skips creating when only difference is case
	db.Model(&data.Tag{}).Where("name = ?", "toRead").Count(&count)
	require.Equal(t, int64(1), count)
	db.Model(&data.Tag{}).Where("name = ?", "TOREAD").Count(&count)
	require.Equal(t, int64(0), count)
	tags, err = repo.Upsert([]string{"TOREAD"})
	require.NoError(t, err)
	require.Equal(t, "toRead", tags[0].Name)
	db.Model(&data.Tag{}).Where("name = ?", "TOREAD").Count(&count)
	require.Equal(t, int64(0), count)
}
