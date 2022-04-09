package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Printf("====Before transfer-> account1: %d : account2: %d\n", account1.Balance, account2.Balance)

	// testing 5 concurrent transactions where each transfers 10(amount) from account1 to account 2
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		fmt.Println("tx ", i)
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: int64(account1.ID),
				ToAccountID:   int64(account2.ID),
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//test transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, int64(account1.ID), transfer.FromAccountID)
		require.Equal(t, int64(account2.ID), transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetAccount(context.Background(), int32(transfer.ID))
		require.NoError(t, err)

		//test from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, int64(account1.ID), fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//test to Entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, int64(account2.ID), toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check accounts are valid
		require.NotEmpty(t, result.FromAccount)
		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, account1.ID, result.FromAccount.ID)
		require.Equal(t, account2.ID, result.ToAccount.ID)

		//check accounts balance
		fmt.Printf("====Inbetween transfer-> account1: %d : account2: %d\n", result.FromAccount.Balance, result.ToAccount.Balance)
		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	//check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, updatedAccount1.Balance+int64(n)*amount, account1.Balance)
	require.Equal(t, updatedAccount2.Balance-int64(n)*amount, account2.Balance)

	fmt.Printf("====Completed transfer-> account1: %d : account2: %d\n", updatedAccount1.Balance, account1.Balance)
	fmt.Printf("====Completed transfer-> account1: %d : account2: %d\n", updatedAccount2.Balance, account2.Balance)

}
