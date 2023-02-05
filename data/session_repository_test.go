package data_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/test"
)

func TestGetByToken(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewSessionRepository(db)

	guid, err := uuid.NewUUID()
	require.NoError(t, err)

	session, err := repo.GetByToken(guid.String())
	require.Nil(t, session)
	require.Nil(t, err)

	result := db.Create(&data.Session{
		Token:     guid.String(),
		UserAgent: "test-agent",
		IP:        "test-ip",
	})
	require.NoError(t, result.Error)

	session, err = repo.GetByToken(guid.String())
	require.Nil(t, err)
	require.NotNil(t, session)
	require.Equal(t, guid.String(), session.Token)
	require.Equal(t, "test-agent", session.UserAgent)
	require.Equal(t, "test-ip", session.IP)
}

func TestCreate(t *testing.T) {
	db, cleanup := test.InitTestDB(t)
	defer cleanup()
	repo := data.NewSessionRepository(db)

	session, err := repo.Create(&data.SessionCreate{
		UserAgent: "test-agent",
		IP:        "test-ip",
	})
	require.Nil(t, err)
	require.NotNil(t, session)
	require.NotEmpty(t, session.Token)
	require.Equal(t, "test-agent", session.UserAgent)
	require.Equal(t, "test-ip", session.IP)
}
