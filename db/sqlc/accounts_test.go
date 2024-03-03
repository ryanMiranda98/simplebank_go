package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ryanMiranda98/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
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

	t.Cleanup(cleanUpDB)
}

func TestGetAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	fetchedAccount, err := testQueries.GetAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)
	require.Equal(t, fetchedAccount.Owner, createdAccount.Owner)
	require.Equal(t, fetchedAccount.Balance, createdAccount.Balance)
	require.Equal(t, fetchedAccount.Currency, createdAccount.Currency)

	t.Cleanup(cleanUpDB)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 3; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

	t.Cleanup(cleanUpDB)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	currentBalance := account.Balance

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: currentBalance + 100,
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, account.ID, updatedAccount.ID)
	require.Equal(t, currentBalance+100, updatedAccount.Balance)
	require.WithinDuration(t, account.CreatedAt, updatedAccount.CreatedAt, time.Second)

	t.Cleanup(cleanUpDB)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Errorf(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

	t.Cleanup(cleanUpDB)
}
