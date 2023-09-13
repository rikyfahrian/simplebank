package db

import (
	"context"
	"techschool/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	CreateRandUser(t)

}

func TestGetUser(t *testing.T) {

	user := CreateRandUser(t)

	row, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, row)

	require.Equal(t, row.Username, user.Username)

}

func CreateRandUser(t *testing.T) User {

	hashPW, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: hashPW,
		FullName:       util.RandomName(),
		Email:          util.RandEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user

}
