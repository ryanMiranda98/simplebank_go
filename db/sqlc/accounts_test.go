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
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
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

func cleanUpDB() {
	testQueries.db.ExecContext(context.Background(), "DELETE FROM accounts;")
	testQueries.db.ExecContext(context.Background(), "DELETE FROM entries;")
	testQueries.db.ExecContext(context.Background(), "DELETE FROM transfers;")
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
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	account3 := createRandomAccount(t)

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)

	require.Len(t, accounts, 3)
	require.Equal(t, accounts[0].Owner, account1.Owner)
	require.Equal(t, accounts[1].Owner, account2.Owner)
	require.Equal(t, accounts[2].Owner, account3.Owner)

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
