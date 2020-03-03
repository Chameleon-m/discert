package model

import "testing"

func TestAccount(t *testing.T) *Account {
	t.Helper()

	return &Account{
		Email:    "user@example.org",
		Password: "password",
	}
}
