package db

import (
	"context"
	"testing"

	"github.com/ryanMiranda98/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPwd, err := util.HashPassword(util.RandomString(10))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username: util.RandomOwner(),
		HashedPassword: hashedPwd,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)

	t.Cleanup(cleanUpDB)
}

func TestGetUser(t *testing.T) {
	createdUser := createRandomUser(t)

	fetchedUser, err := testQueries.GetUser(context.Background(), createdUser.Username)
	require.NoError(t, err)
	require.Equal(t, fetchedUser.Email, createdUser.Email)
	require.Equal(t, fetchedUser.FullName, createdUser.FullName)
	require.Equal(t, fetchedUser.Username, createdUser.Username)

	t.Cleanup(cleanUpDB)
}
