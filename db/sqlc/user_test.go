package db

import (
	"context"
	"fmt"
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

	arg := CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: "secret",
		FullName:       util.RandomName(),
		Email:          RandEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user

}

func RandEmail() string {
	return fmt.Sprintf("%s@gmail.com", util.RandomString(6))
}
