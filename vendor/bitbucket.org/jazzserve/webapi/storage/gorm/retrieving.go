package gorm

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"reflect"
)

func (s *Storage) GetAnyOf(modelPtr interface{}) error {
	t := reflect.TypeOf(modelPtr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("getAnyOf expects a pointer to a struct")
	}

	err := s.Find(modelPtr).Error
	if err == gorm.ErrRecordNotFound {
		return err
	}

	return errors.WithStack(err)
}

func (s *Storage) GetAnyOfOptional(modelPtr interface{}) (bool, error) {
	err := s.GetAnyOf(modelPtr)
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return true, errors.WithStack(err)
}

func (s *Storage) GetListOf(arrPtr interface{}) error {
	t := reflect.TypeOf(arrPtr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice {
		return errors.New("getListOf expects a pointer to a slice")
	}

	return errors.WithStack(s.Find(arrPtr).Error)
}

func (s *Storage) FillListOf(arrPtr interface{}) (bool, error) {
	t, v := reflect.TypeOf(arrPtr), reflect.ValueOf(arrPtr)
	if t.Kind() != reflect.Ptr {
		return false, errors.New("fillListOf expects a pointer to a slice")
	}

	t, v = t.Elem(), v.Elem()
	if t.Kind() != reflect.Slice {
		return false, errors.New("fillListOf expects a pointer to a slice")
	}

	for i := 0; i < v.Len(); i++ {
		elemV := v.Index(i)

		if !elemV.IsValid() || !elemV.CanInterface() {
			return false, errors.New("can't fill list of: invalid elem type")
		}

		ok, err := s.GetAnyOfOptional(elemV.Interface())
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}
