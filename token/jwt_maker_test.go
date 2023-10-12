package token

import (
	"techschool/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {

	maker, err := NewJWTMaker(util.RandomString(40))
	require.NoError(t, err)

	username := util.RandomName()

	token, _, err := maker.CreateToken(username, time.Minute)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)

	require.Equal(t, payload.Username, username)

}
