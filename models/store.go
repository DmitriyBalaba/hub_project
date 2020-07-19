package models

import (
	"fmt"

	"bitbucket.org/jazzserve/webapi/utils/pointers/val"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Store struct {
	ID               *int64          `json:"id,omitempty"`
	Name             *string         `json:"name,omitempty"`
	AccountStoreList []*AccountStore `json:"account_store_list,omitempty"`
}

func (s *storage) UpdateStore(t *Store) error {
	fnTx := func(tx *gorm.DB) error {
		if err := tx.Model(&Store{}).Where("id = ?", t.ID).Update(&t).Error; err != nil {
			return err
		}

		for _, a := range t.AccountStoreList {
			upAccountStore := &AccountStore{
				StoreID:   t.ID,
				AccountID: a.ID,
			}
			err := tx.Assign(&upAccountStore).FirstOrCreate(&a, val.Int64(a.ID)).Error
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}

	return s.RunInTransaction(fnTx)
}

func (s *storage) DeleteStore(ID *int64) error {
	if ID == nil {
		return fmt.Errorf("stroe's id is reauired")
	}
	return s.DB.Delete(&Store{ID: ID}).Association("AccountList").Clear().Error
}
