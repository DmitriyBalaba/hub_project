package models

import (
	"fmt"

	psql "bitbucket.org/jazzserve/webapi/storage/gorm"
	"bitbucket.org/jazzserve/webapi/utils"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Storage interface {
	Match(field string, value interface{}) Storage
	SortAsc(field string) Storage
	SortDesc(field string) Storage
	OffsetAndLimit(offset, limit int64) Storage
	GetListOf(arrPtr interface{}) error
	GetAnyOf(arrPtr interface{}) error
	GetAnyOfOptional(arrPtr interface{}) (bool, error)
	Preload(column string, conditions ...interface{}) Storage
	LeftJoin(query string, args ...interface{}) Storage
	MatchEmail(email string) Storage

	//CRUD
	Update(model interface{}) error
	Create(model interface{}) error
	Delete(model interface{}) error

	//migrate.go
	Migrate(auto []interface{}, manual []*Migration) error

	//store.go
	UpdateStore(t *Store) error
	DeleteStore(ID *int64) error

	//accountStore.go
	CreateAccountStore(as *AccountStore) error
	PreloadAccount(accountID *int64) Storage
}

const (
	DefaultOffset = 0
	DefaultLimit  = 30
)

type PrepareQuerier interface {
	PrepareQuery(s Storage) Storage
}

type Tabler interface {
	Table() string
}

type Updater interface {
	Update(s Storage) error
}

type storage struct {
	*psql.Storage
}

func NewStorage(db *gorm.DB) Storage {
	s := &storage{
		Storage: psql.NewStorage(db),
	}
	return wrap(s.SetDefaultOffsetAndLimit(DefaultOffset, DefaultLimit))
}

func wrap(store *psql.Storage) Storage {
	return &storage{store}
}

func (s *storage) Match(field string, value interface{}) Storage {
	if utils.IsZeroAny(value) {
		return s
	}
	return wrap(s.Storage.Match(field, value))
}

func (s *storage) GetListOf(arrPtr interface{}) error {
	return errors.WithStack(s.Storage.GetListOf(arrPtr))
}

func (s *storage) GetAnyOf(arrPtr interface{}) error {
	return errors.WithStack(s.Storage.GetAnyOf(arrPtr))
}

func (s *storage) GetAnyOfOptional(arrPtr interface{}) (bool, error) {
	exists, err := s.Storage.GetAnyOfOptional(arrPtr)
	if err != nil {
		err = errors.WithStack(err)
	}
	return exists, err
}

func (s *storage) Update(model interface{}) error {
	return errors.WithStack(s.Model(model).Update(model).Error)
}

func (s *storage) OffsetAndLimit(offset, limit int64) Storage {
	return wrap(s.Storage.OffsetAndLimit(offset, limit))
}

func (s *storage) SortAsc(field string) Storage {
	return wrap(s.Storage.SortAsc(field))
}

func (s *storage) SortDesc(field string) Storage {
	return wrap(s.Storage.SortDesc(field))
}

func (s *storage) Delete(model interface{}) error {
	return psql.Wrap(s.Storage.Delete(model)).Error
}

func (s *storage) Create(model interface{}) error {
	return errors.WithStack(s.Storage.Create(model).Error)
}

func (s *storage) LeftJoin(query string, args ...interface{}) Storage {
	return wrap(psql.Wrap(s.Storage.Joins(fmt.Sprintf("LEFT JOIN %s", query), args)))
}

func (s *storage) RunInTransaction(operations ...psql.DbOperation) error {
	return errors.WithStack(s.Storage.RunInTransaction(operations...))
}

func (s *storage) MatchEmail(email string) Storage {
	email = unifyStr(email)
	if email == "" {
		return s
	}
	return s.matchNoCase("email", email)
}

func (s *storage) matchNoCase(field, value string) Storage {
	field, value = unifyStr(field), unifyStr(value)
	if utils.IsZeroAny(field, value) {
		return s
	}
	query := fmt.Sprintf("lower(%s) = ?", field)
	return NewStorage(s.Storage.Where(query, value))
}

func (s *storage) Preload(column string, conditions ...interface{}) Storage {
	return NewStorage(s.Storage.Preload(column, conditions))
}
