package gorm

import (
	"bitbucket.org/jazzserve/webapi/utils/pointers/val"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Storage struct {
	*gorm.DB
}

const (
	defaultOffsetKey = "default-offset"
	defaultOffsetVal = 0
	defaultLimitKey  = "default-limit"
	defaultLimitVal  = 10
)

func Wrap(db *gorm.DB) *Storage {
	return &Storage{db}
}

func NewStorage(db *gorm.DB) *Storage {
	return Wrap(db.
		Set(defaultOffsetKey, defaultOffsetVal).
		Set(defaultLimitKey, defaultLimitVal))
}

func (s *Storage) GetStrByKey(key string) (string, error) {
	str, ok := s.Get(key)
	if !ok {
		return "", errors.New("key not found: " + key)
	}

	if v, ok := str.(string); ok {
		return v, nil
	}

	if vPtr, ok := str.(*string); ok {
		return val.Str(vPtr), nil
	}

	return "", errors.Errorf("%s has undefined type: expected string", key)
}

func (s *Storage) GetInt64ByKey(key string) (int64, error) {
	i, ok := s.Get(key)
	if !ok {
		return 0, errors.New("key not found: " + key)
	}

	if v, ok := i.(int64); ok {
		return v, nil
	}

	if vPtr, ok := i.(*int64); ok {
		return val.Int64(vPtr), nil
	}

	return 0, errors.Errorf("%s has undefined type: expected int", key)
}

func (s *Storage) SetDefaultOffsetAndLimit(offset, limit int64) *Storage {
	return Wrap(s.
		Set(defaultOffsetKey, offset).
		Set(defaultLimitKey, limit))
}

func (s *Storage) GetDefaultOffsetAndLimit() (offset int64, limit int64) {
	offset, err := s.GetInt64ByKey(defaultOffsetKey)
	if err != nil {
		panic("can't retrieve default offset")
	}

	limit, err = s.GetInt64ByKey(defaultLimitKey)
	if err != nil {
		panic("can't retrieve default limit")
	}

	return
}

type DbOperation func(db *gorm.DB) error

func (s *Storage) RunInTransaction(operations ...DbOperation) (err error) {
	tx := s.Begin()
	defer func(txCp *gorm.DB) {
		if r := recover(); r != nil {
			txCp.Rollback()
			err = errors.Errorf("%v", r)
		}
	}(tx)

	if err = tx.Error; err != nil {
		return
	}

	for _, f := range operations {
		if err = f(tx); err != nil {
			tx.Rollback()
			return
		}
	}

	return tx.Commit().Error
}
