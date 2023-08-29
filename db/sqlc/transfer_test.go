package db

import (
	"context"
	"techschool/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateTransferForTest(t *testing.T) Transfer {

	account1 := CreateAccountForTest(t)
	account2 := CreateAccountForTest(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}

	result, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)

	return result

}

func TestCreateTransfer(t *testing.T) {
	CreateTransferForTest(t)

}

func TestGetTransfers(t *testing.T) {
	tf := CreateTransferForTest(t)

	result, err := testQueries.GetTransfer(context.Background(), tf.ID)

	require.NoError(t, err)
	require.Equal(t, result.ID, tf.ID)

}
