package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/kizitonzeka/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	args := CreateEntryParams{
		AccountID: int64(account.ID),
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)

	return entry
}
func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestUpdateEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	args := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}

	err := testQueries.UpdateEntry(context.Background(), args)
	require.NoError(t, err)
	entry2, err := testQueries.GetEntry(context.Background(), args.ID)
	require.NoError(t, err)

	require.NotEmpty(t, entry2)
	require.Equal(t, args.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, args.Amount, entry2.Amount)

}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
}

func TestListEntries(t *testing.T) {
	for i := 1; i <= 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entryItems, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entryItems, 5)

	for _, entry := range entryItems {
		require.NotEmpty(t, entry)
	}

}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)

}
