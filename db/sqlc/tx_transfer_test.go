package db

import (
	"context"
	// "fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	// fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transaction
	n := 5
	amount := int64(10)

	// copnfiguring two channels> one for error and another for transfertxresult to get result from concurrent functions.
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check Entries
		FromEntry := result.FromEntry
		require.NotEmpty(t, FromEntry)
		require.Equal(t, account1.ID, FromEntry.AccountID)
		require.Equal(t, -amount, FromEntry.Amount)
		require.NotZero(t, FromEntry.ID)
		require.NotZero(t, FromEntry.CreatedAt)
		_, err = testStore.GetEntry(context.Background(), FromEntry.ID)
		require.NoError(t, err)

		ToEntry := result.ToEntry
		require.NotEmpty(t, ToEntry)
		require.Equal(t, account2.ID, ToEntry.AccountID)
		require.Equal(t, amount, ToEntry.Amount)
		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)
		_, err = testStore.GetEntry(context.Background(), ToEntry.ID)
		require.NoError(t, err)

		// check accounts
		FromAccount := result.FromAccount
		require.NotEmpty(t, FromAccount)
		require.Equal(t, account1.ID, FromAccount.ID)

		ToAccount := result.ToAccount
		require.NotEmpty(t, ToAccount)
		require.Equal(t, account2.ID, ToAccount.ID)

		// check accounts' balance
		// fmt.Println(">> tx:", FromAccount.Balance, ToAccount.Balance)

		diff1 := account1.Balance - FromAccount.Balance
		diff2 := ToAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

// testing deadlock when transactions are done back and forth from account 1 and 2
func TestTransferTxDeadlock(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	// fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transaction
	n := 10
	amount := int64(10)

	// configuring one channel only because results have already been checked
	errs := make(chan error)

	for i := 0; i < n; i++ {

		go func() {
			fromAccountID := account1.ID
			toAccountID := account2.ID

			// half transactions with 1>2 and the other half from 2>1
			if i%2 == 1 {
				fromAccountID = account2.ID
				toAccountID = account1.ID
			}

			ctx := context.Background()
			_, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balances
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
