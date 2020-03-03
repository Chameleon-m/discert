package teststore_test

import (
	"github.com/Chameleon-m/discert/internal/app/model"
	"github.com/Chameleon-m/discert/internal/app/store"
	"github.com/Chameleon-m/discert/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountRepository_Create(t *testing.T) {

	s := teststore.New()
	a := model.TestAccount(t)
	err := s.AccountRepository().Create(a)

	assert.NoError(t, err)
	assert.NotNil(t, a)
}

func TestAccountRepository_Find(t *testing.T) {
	s := teststore.New()
	id := int64(1)
	_, err := s.AccountRepository().Find(id)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	account1 := model.TestAccount(t)
	account1.ID = id
	s.AccountRepository().Create(account1)
	account2, err := s.AccountRepository().Find(account1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, account2)
}

func TestAccountRepository_FindByEmail(t *testing.T) {

	s := teststore.New()
	email := "user@example.org"
	_, err := s.AccountRepository().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	a := model.TestAccount(t)
	a.Email = email
	s.AccountRepository().Create(a)
	a, err = s.AccountRepository().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, a)
}
