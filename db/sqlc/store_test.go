package db

import (
	"context"
	"fmt"
	"sync"

	"testing"

	"gopkg.in/stretchr/testify.v1/require"
)

func TestTransferTx(t *testing.T) {
	Store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> Before:", account1.Balance, account2.Balance)

	// Run n concurrent transfer transactions
	n := 5
	amount := int64(10)
	resultChan := make(chan TransferResult, n)
	errChan := make(chan error, n)
	var wg sync.WaitGroup

	wg.Add(n)
	for i := 0; i < n; i++ {
		// txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			// ctx := context.WithValue(context.Background(), txKey, txName)
			ctx := context.Background()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("recovered from", r)
					errChan <- nil
					resultChan <- TransferResult{}
					wg.Done()
				}
			}()

			result, err := Store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errChan <- err
			resultChan <- result
			wg.Done()
		}()
	}

	wg.Wait()

	// Check for errors
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)

		result := <-resultChan
		require.NotEmpty(t, result)

		// check transfer status\
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = Store.GetTranfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// Check entries
		FromEntry := result.FromEntry
		require.NotEmpty(t, FromEntry)
		require.Equal(t, FromEntry.AccountID, account1.ID)
		require.Equal(t, FromEntry.Amount, -amount)
		require.NotZero(t, FromEntry.ID)
		require.NotZero(t, FromEntry.CreatedAt)

		_, err = Store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		// Check entries
		ToEntry := result.ToEntry
		require.NotEmpty(t, ToEntry)
		require.Equal(t, ToEntry.AccountID, account2.ID)
		require.Equal(t, ToEntry.Amount, amount)
		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)

		_, err = Store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		// check accounts balance
		fmt.Println(">>Tx: ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // Check only with case amount = 10 has been declared to test

		k := int(diff1 / amount)
		require.True(t, k > 0 && k <= n)
		existed[k] = true
	}

	// check the final update balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>After: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	Store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> Before:", account1.Balance, account2.Balance)
	n := 10
	amount := int64(10)
	errChan := make(chan error, n)
	var wg sync.WaitGroup

	wg.Add(n)
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("recovered from", r)
					errChan <- nil
					wg.Done()
				}
			}()
			_, err := Store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errChan <- err
			wg.Done()
		}()
	}

	wg.Wait()

	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)
	}

	// check the final update balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>After: ", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
