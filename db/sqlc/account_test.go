package db

import (
	"context"
	"database/sql"
	"simplebank/utils"
	"testing"
	"time"

	"gopkg.in/stretchr/testify.v1/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMonney(),
		Currency: utils.RandomCurrencies(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	newAccount := createRandomAccount(t)
	getAccount, err := testQueries.GetAccount(context.Background(), newAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getAccount)
	require.Equal(t, newAccount.ID, getAccount.ID)
	require.Equal(t, newAccount.Owner, getAccount.Owner)
	require.Equal(t, newAccount.Balance, getAccount.Balance)
	require.Equal(t, newAccount.Currency, getAccount.Currency)
	require.WithinDuration(t, newAccount.CreatedAt, getAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	newAccount := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      newAccount.ID,
		Balance: utils.RandomMonney(),
	}
	updateAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount)
	require.Equal(t, newAccount.ID, updateAccount.ID)
	require.Equal(t, newAccount.Owner, updateAccount.Owner)
	require.Equal(t, arg.Balance, updateAccount.Balance)
	require.Equal(t, newAccount.Currency, updateAccount.Currency)
	require.WithinDuration(t, newAccount.CreatedAt, updateAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	newAccount := createRandomAccount(t)
    err := testQueries.DeleteAccount(context.Background(), newAccount.ID)
    require.NoError(t, err)

	getAccount, err := testQueries.GetAccount(context.Background(), newAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, getAccount)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

