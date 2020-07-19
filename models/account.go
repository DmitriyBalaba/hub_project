package models

import (
	"errors"
	"fmt"

	"bitbucket.org/jazzserve/webapi/utils"
	"bitbucket.org/jazzserve/webapi/utils/pointers/ptr"
	"bitbucket.org/jazzserve/webapi/utils/pointers/val"
	"github.com/jinzhu/gorm"
)

type Account struct {
	ID             *int64        `json:"id,omitempty"`
	Name           *string       `json:"name,omitempty"`
	Email          *string       `json:"email,omitempty"`
	IsAdmin        *bool         `json:"is_admin,omitempty" gorm:"default:false"`
	IsStoreManager *bool         `json:"is_store_manager,omitempty" gorm:"default:false"`
	Password       *string       `json:"password,omitempty"`
	AccountStore   *AccountStore `json:"account_store,omitempty"`
}

func (a *Account) HasEmail(e string) bool {
	return unifyStr(val.Str(a.Email)) == unifyStr(e)
}

func (a *Account) HasPassword(p string) bool {
	return val.Str(a.Password) == p
}

func IsAccountExists(email *string, ID *int64, s Storage) (bool, error) {
	if val.Str(email) == "" {
		return true, nil
	}

	acc := &Account{}
	err := s.Match("email", val.Str(email)).GetAnyOf(acc)
	if errors.Is(err, gorm.ErrRecordNotFound) || val.Int64(ID) == val.Int64(acc.ID) {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return false, nil
}

func (a *Account) Create(s Storage) error {
	if a == nil || val.Str(a.Email) == "" {
		return fmt.Errorf("email is required")
	}

	password, err := utils.GenerateRandomString(8)
	if err != nil {
		return err
	}

	a.Password = ptr.Str(password)
	return s.Create(&a)
}
