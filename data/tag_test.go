package data_test

import (
	"testing"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/test"
	"github.com/stretchr/testify/require"
)

func TestTagBeforeSave(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()

	tag := data.Tag{DisplayName: "ActivityPub"}
	db.Create(&tag)

	require.Equal(t, "ActivityPub", tag.DisplayName)
	require.Equal(t, "activitypub", tag.Name)
}
