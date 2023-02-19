package data_test

import (
	"testing"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/test"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()

	migrator := data.NewMigrator(db)

	for {
		err := migrator.RollbackLast()
		if err == gormigrate.ErrNoRunMigration {
			break
		}
		require.NoError(t, err)
	}
	err := migrator.Migrate()
	require.NoError(t, err)
}
