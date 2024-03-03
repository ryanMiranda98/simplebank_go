package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func generateHashedPassword(t *testing.T) string {
	password := "Secret1234"
	hashedPwd, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEqual(t, password, hashedPwd)

	return hashedPwd
}

func TestHashPassword(t *testing.T) {
	generateHashedPassword(t)
}

func TestHashPasswordError(t *testing.T) {
	password := ""
	hashedPwd, err := HashPassword(password)
	require.Empty(t, hashedPwd)
	require.ErrorContains(t, err, "invalid password provided")

	longPassword := RandomString(75)
	hashedPwd, err = HashPassword(longPassword)
	require.Empty(t, hashedPwd)
	require.ErrorContains(t, err, "failed to hash password")
}

func TestComparePassword(t *testing.T) {
	hashedPwd := generateHashedPassword(t)
	err := ComparePassword(hashedPwd, "Secret1234")
	require.NoError(t, err)
}

func TestComparePasswordError(t *testing.T) {
	hashedPwd := generateHashedPassword(t)
	err := ComparePassword(hashedPwd, "Secret14")
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}