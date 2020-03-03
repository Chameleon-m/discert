package model_test

import (
	"github.com/Chameleon-m/discert/internal/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccount_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		account func() *model.Account
		isValid bool
	}{
		{
			name: "valid",
			account: func() *model.Account {
				return model.TestAccount(t)
			},
			isValid: true,
		},
		{
			name: "empty email",
			account: func() *model.Account {
				a := model.TestAccount(t)
				a.Email = ""

				return a
			},
			isValid: false,
		},
		{
			name: "invalid email",
			account: func() *model.Account {
				a := model.TestAccount(t)
				a.Email = "invalid"

				return a
			},
			isValid: false,
		},
		{
			name: "empty password",
			account: func() *model.Account {
				a := model.TestAccount(t)
				a.Password = ""

				return a
			},
			isValid: false,
		},
		{
			name: "short password",
			account: func() *model.Account {
				a := model.TestAccount(t)
				a.Password = "short"

				return a
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.account().Validate())
			} else {
				assert.Error(t, tc.account().Validate())
			}
		})
	}
}

func TestAccount_BeforeCreate(t *testing.T) {
	a := model.TestAccount(t)
	assert.NoError(t, a.BeforeCreate())
	assert.NotEmpty(t, a.Password)
}
