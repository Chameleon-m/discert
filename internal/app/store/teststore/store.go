package teststore

import (
	"github.com/Chameleon-m/discert/internal/app/model"
	"github.com/Chameleon-m/discert/internal/app/store"
)

type Store struct {
	accountRepository *AccountRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) AccountRepository() store.AccountRepository {
	if s.accountRepository != nil {
		return s.accountRepository
	}

	s.accountRepository = &AccountRepository{
		store:    s,
		accounts: make(map[int64]*model.Account),
	}

	return s.accountRepository
}
