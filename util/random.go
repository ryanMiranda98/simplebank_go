package util

import (
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	alphabets string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// RandomInt generates a random integer between min and max values
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Random string generates a random string with n length
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency
func RandomCurrency() string {
	currencies := []string{"USD", "INR", "CAD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	email := fmt.Sprintf("%s@email.com", RandomString(6))
	return email
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("invalid password provided")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
