package test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/chdorner/submarine/data"
	"github.com/stretchr/testify/require"
)

func InitTestDB(t *testing.T) (*gorm.DB, func()) {
	path := fmt.Sprintf("./test-%d.db", rand.Int())

	db, err := data.Connect(path)
	require.NoError(t, err)
	data.Migrate(db)

	return db, func() {
		os.Remove(path)
		os.Remove(fmt.Sprintf("%s-journal", path))
	}
}
