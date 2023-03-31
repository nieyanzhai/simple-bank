package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/nieyanzhai/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestExecTransferTx(t *testing.T) {
	// Create a new store
	store := NewStore(testDB)

	// Create two accounts
	account1, err := store.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandString(6),
		Balance:  100,
		Currency: util.RandCurrency(),
	})
	require.NoError(t, err)
	account2, err := store.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandString(6),
		Balance:  0,
		Currency: util.RandCurrency(),
	})
	require.NoError(t, err)

	fmt.Printf(">> before transfer: %d, %d\n", account1.Balance, account2.Balance)

	// Transfer amount of money from account1 to account2 for n times
	n := 10
	amount := int64(10)
	errCh := make(chan error)
	resultCh := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.ExecTransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errCh <- err
			resultCh <- result
		}()
	}

	// Check if all the transfers were successful
	for i := 0; i < n; i++ {
		err := <-errCh
		require.NoError(t, err)

		// Check if the transfer was successful
		result := <-resultCh
		require.NotEmpty(t, result)

		require.NotEmpty(t, result.Transfer)
		require.NotZero(t, result.Transfer.ID)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.CreatedAt)

		require.NotEmpty(t, result.FromEntry)
		require.NotZero(t, result.FromEntry.ID)
		require.Equal(t, account1.ID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.CreatedAt)

		require.NotEmpty(t, result.ToEntry)
		require.NotZero(t, result.ToEntry.ID)
		require.Equal(t, account2.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.CreatedAt)

		fmt.Printf(">> after transfer: %d, %d\n", result.FromAccount.Balance, result.ToAccount.Balance)

		// Check if the account balances were updated correctly
		diff1 := result.FromAccount.Balance - account1.Balance
		require.Equal(t, diff1, -amount*int64(i+1))
		diff2 := result.ToAccount.Balance - account2.Balance
		require.Equal(t, diff2, amount*int64(i+1))
	}

	// Check if the accounts' balances were updated correctly
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	diff := updatedAccount1.Balance - account1.Balance
	require.NoError(t, err)
	require.Equal(t, diff, -amount*int64(n))

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	diff = updatedAccount2.Balance - account2.Balance
	require.NoError(t, err)
	require.Equal(t, diff, amount*int64(n))
}

func TestExecTransferTx_DeadLock(t *testing.T) {
	// Create a new store
	store := NewStore(testDB)

	// Create two accounts
	account1, err := store.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandString(6),
		Balance:  100,
		Currency: util.RandCurrency(),
	})
	require.NoError(t, err)
	account2, err := store.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandString(6),
		Balance:  100,
		Currency: util.RandCurrency(),
	})
	require.NoError(t, err)

	fmt.Printf(">> before transfer: %d, %d\n", account1.Balance, account2.Balance)

	// Transfer amount of money from account1 to account2 for n times
	n := 10
	amount := int64(10)
	errCh := make(chan error)
	for i := 0; i < n; i++ {
		var fromAccountID int64
		var toAccountID int64
		if i%2 == 1 {
			fromAccountID = account1.ID
			toAccountID = account2.ID
		} else {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.ExecTransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errCh <- err
		}()
	}

	// Check if all the transfers were successful
	for i := 0; i < n; i++ {
		err := <-errCh
		require.NoError(t, err)
	}

	// Check if the accounts' balances were updated correctly
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	diff := updatedAccount1.Balance - account1.Balance
	require.NoError(t, err)
	require.Equal(t, diff, int64(0))

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	diff = updatedAccount2.Balance - account2.Balance
	require.NoError(t, err)
	require.Equal(t, diff, int64(0))

	fmt.Printf(">> after transfer: %d, %d\n", updatedAccount1.Balance, updatedAccount2.Balance)
}
