package test

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/chdorner/submarine/data"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func InitTestDB(t *testing.T) (*gorm.DB, func()) {
	guid, err := uuid.NewUUID()
	require.NoError(t, err)

	path := fmt.Sprintf("./test-%s.db", guid.String())

	db, err := data.Connect(path)
	require.NoError(t, err)
	data.Migrate(db)

	return db, func() {
		os.Remove(path)
		os.Remove(fmt.Sprintf("%s-journal", path))
	}
}
