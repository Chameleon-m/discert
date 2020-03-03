package sqlstore

import (
	"database/sql"
	"github.com/Chameleon-m/discert/internal/app/store"
	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db                *sql.DB
	accountRepository *AccountRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) AccountRepository() store.AccountRepository {
	if store.accountRepository != nil {
		return store.accountRepository
	}

	store.accountRepository = &AccountRepository{
		store: store,
	}

	return store.accountRepository
}
