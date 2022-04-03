package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kizitonzeka/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: int64(account1.ID),
		ToAccountID:   int64(account2.ID),
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.CreatedAt)

	return transfer

}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)

	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)

}

func TestListTransfers(t *testing.T) {
	for i := 1; i <= 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestUpdateTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomMoney(),
	}

	err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	transfer2, err := testQueries.GetTransfer(context.Background(), arg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, arg.ID, transfer2.ID)
	require.Equal(t, arg.Amount, transfer2.Amount)

}

func TestDeleteTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}
