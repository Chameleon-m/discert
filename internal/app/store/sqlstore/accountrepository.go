package sqlstore

import (
	"database/sql"
	"github.com/Chameleon-m/discert/internal/app/model"
	"github.com/Chameleon-m/discert/internal/app/store"
)

type AccountRepository struct {
	store *Store
}

func (r *AccountRepository) Create(a *model.Account) error {
	if err := a.Validate(); err != nil {
		return err
	}

	if err := a.BeforeCreate(); err != nil {
		return err
	}

	result, err := r.store.db.Exec(
		"INSERT INTO accounts (email, phone, password) VALUES (?,?,?)",
		a.Email,
		a.Phone,
		a.Password,
	)
	if err != nil {
		return err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = lastInsertId

	return nil
}

func (r *AccountRepository) Find(id int64) (*model.Account, error) {
	a := &model.Account{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, phone, password FROM accounts WHERE id = ?",
		id,
	).Scan(
		&a.ID,
		&a.Email,
		&a.Phone,
		&a.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return a, nil
}

func (r *AccountRepository) FindByEmail(email string) (*model.Account, error) {
	a := &model.Account{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, phone, password FROM accounts WHERE email = ?",
		email,
	).Scan(
		&a.ID,
		&a.Email,
		&a.Phone,
		&a.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return a, nil
}
