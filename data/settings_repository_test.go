package data_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/test"
)

func TestIsInitialized(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewSettingsRepository(db)

	require.False(t, repo.IsInitialized())

	err := repo.Upsert(data.SettingsUpsert{Password: "supersecret"})
	require.NoError(t, err)

	require.True(t, repo.IsInitialized())
}

func TestSettingsRepositoryGet(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewSettingsRepository(db)

	actual, err := repo.Get()
	require.Nil(t, actual)
	require.NoError(t, err)

	err = repo.Upsert(data.SettingsUpsert{Password: "supersecret"})
	require.NoError(t, err)

	actual, err = repo.Get()
	require.NoError(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(actual.Password), []byte("supersecret"))
	require.NoError(t, err)
}

func TestSettingsRepositoryUpsert(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewSettingsRepository(db)

	// insert when settings table is empty
	var count int64
	db.Model(&data.Settings{}).Count(&count)
	require.Equal(t, int64(0), count)

	// insert fails when password is empty
	err := repo.Upsert(data.SettingsUpsert{Password: ""})
	require.EqualError(t, err, "password is empty")

	// insert succeeds
	err = repo.Upsert(data.SettingsUpsert{Password: "secret"})
	require.NoError(t, err)
	var actual data.Settings
	db.First(&actual)
	err = bcrypt.CompareHashAndPassword([]byte(actual.Password), []byte("secret"))
	require.NoError(t, err)

	// update password
	err = repo.Upsert(data.SettingsUpsert{Password: "topsecret"})
	require.NoError(t, err)
	db.First(&actual)
	err = bcrypt.CompareHashAndPassword([]byte(actual.Password), []byte("topsecret"))
	require.NoError(t, err)
}
