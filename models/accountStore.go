package models

import "github.com/jinzhu/gorm"

type AccountStore struct {
	ID        *int64 `json:"id,omitempty"`
	AccountID *int64 `json:"-"`
	StoreID   *int64 `json:"-"`

	Account *Account `json:"account,omitempty" gorm:"association_autoupdate:false"`
	Store   *Store   `json:"store,omitempty" gorm:"association_autoupdate:false"`
}

func (s *AccountStore) AutoMigrate(tx *gorm.DB) error {
	return tx.Model(&AccountStore{}).
		AddForeignKey("account_id", "account(id)", "CASCADE", "CASCADE").
		AddIndex("idx_account_store_account_id", "account_id").
		AddForeignKey("store_id", "store(id)", "CASCADE", "CASCADE").
		AddIndex("idx_account_store_store_id", "store_id").
		Error
}

func (s *storage) CreateAccountStore(as *AccountStore) error {

	fnTx := func(tx *gorm.DB) error {
		err := tx.Where("account_id = ?", as.Account.ID).Delete(&AccountStore{}).Error
		if err != nil {
			return err
		}

		return tx.Create(&as).Error
	}
	return s.RunInTransaction(fnTx)
}

func (s *storage) PreloadAccount(accountID *int64) Storage {
	return NewStorage(s.Storage.Preload("AccountStoreList", func(tx *gorm.DB) *gorm.DB {
		return tx.Where("account_id = ?", accountID)
	}).Preload("AccountStoreList.Account"))
}
