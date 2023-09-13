package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasswordHash(t *testing.T) {

	pw := RandomString(8)

	pwHash, err := HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, pwHash)

	err = CheckPassword(pw, pwHash)
	require.NoError(t, err)

	wrongPw := RandomString(8)
	err = CheckPassword(wrongPw, pwHash)
	require.Error(t, err)

}
