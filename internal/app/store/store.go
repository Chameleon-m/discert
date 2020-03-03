package store

type Store interface {
	AccountRepository() AccountRepository
}
