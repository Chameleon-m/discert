package store

import "github.com/Chameleon-m/discert/internal/app/model"

type AccountRepository interface {
	Create(*model.Account) error
	FindByEmail(string) (*model.Account, error)
	Find(int64) (*model.Account, error)
}
