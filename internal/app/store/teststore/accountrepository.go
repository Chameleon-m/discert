package teststore

import (
	"github.com/Chameleon-m/discert/internal/app/model"
	"github.com/Chameleon-m/discert/internal/app/store"
)

type AccountRepository struct {
	store    *Store
	accounts map[int64]*model.Account
}

func (r *AccountRepository) Create(account *model.Account) error {
	if err := account.Validate(); err != nil {
		return err
	}

	if err := account.BeforeCreate(); err != nil {
		return err
	}

	account.ID = int64(len(r.accounts)) + 1
	r.accounts[account.ID] = account

	return nil
}

func (r *AccountRepository) Find(id int64) (*model.Account, error) {
	a, ok := r.accounts[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return a, nil
}

func (r *AccountRepository) FindByEmail(email string) (*model.Account, error) {
	for _, account := range r.accounts {
		if account.Email == email {
			return account, nil
		}
	}

	return nil, store.ErrRecordNotFound
}
