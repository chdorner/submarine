package util_test

import (
	"testing"

	"github.com/chdorner/submarine/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	actual, err := util.HashPassword("secret")
	require.NoError(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(actual), []byte("secret"))
	require.NoError(t, err)

	_, err = util.HashPassword("")
	require.NoError(t, err)

	_, err = util.HashPassword("this password is way too long, like really really long, this still isn't long enough, but this is too long")
	require.Error(t, err)
}

func TestComparePassword(t *testing.T) {
	hashed, err := util.HashPassword("supersecret")
	require.NoError(t, err)

	require.True(t, util.ComparePassword("supersecret", hashed))
	require.False(t, util.ComparePassword("secret", hashed))
	require.False(t, util.ComparePassword("this password is way too long, like really really long, this still isn't long enough, but this is too long", hashed))
}
