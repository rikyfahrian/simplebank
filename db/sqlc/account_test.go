package db

import (
	"context"
	"techschool/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateAccountForTest(t *testing.T) Account {

	user := CreateRandUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	result, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, arg.Balance, result.Balance)
	require.Equal(t, arg.Currency, result.Currency)
	require.Equal(t, arg.Owner, result.Owner)

	return result
}

func TestCreateAccount(t *testing.T) {
	CreateAccountForTest(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateAccountForTest(t)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1, account2)
}

func TestUpdateAccount(t *testing.T) {

	account1 := CreateAccountForTest(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomBalance(),
	}

	result, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqual(t, result.Balance, account1.Balance)

}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateAccountForTest(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	result, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.Equal(t, result.ID, int64(0))

}

func TestGetListAccount(t *testing.T) {

	var lastAccount Account

	for i := 0; i < 10; i++ {
		lastAccount = CreateAccountForTest(t)

	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	result, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, result)
}
