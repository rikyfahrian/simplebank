package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(testDB)

	// n := 5

	account1 := CreateAccountForTest(t)
	account2 := CreateAccountForTest(t)

	log.Println("before: ", account1.Balance, account2.Balance)

	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	go func() {

		result, err := store.TransferTx(context.Background(), TransferTxParams{
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        amount,
		})

		errs <- err
		results <- result

	}()

	err := <-errs
	require.NoError(t, err)

	result := <-results
	require.NotEmpty(t, result)

	// check transfer
	transfer := result.Transfer

	log.Println(result.ToAccount.Balance)
	require.NotEmpty(t, transfer)
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

}
