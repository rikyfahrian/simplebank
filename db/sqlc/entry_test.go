package db

import (
	"context"
	"techschool/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateEntryForTest(t *testing.T) Entry {
	account1 := CreateAccountForTest(t)

	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    util.RandomBalance(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, entry.AccountID, account1.ID)

	return entry
}
func TestCreateEntry(t *testing.T) {
	CreateEntryForTest(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := CreateEntryForTest(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	require.Equal(t, entry1.Amount, entry2.Amount)

}
