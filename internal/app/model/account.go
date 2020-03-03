package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Phone    int    `json:"phone"`
	Password string `json:"-"`
}

func (a *Account) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.Required, validation.Length(6, 100)),
	)
}

func (a *Account) BeforeCreate() error {
	if len(a.Password) > 0 {
		enc, err := encryptString(a.Password)
		if err != nil {
			return err
		}

		a.Password = enc
	}
	return nil
}

func (a *Account) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}
