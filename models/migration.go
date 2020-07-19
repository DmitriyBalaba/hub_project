package models

import (
	"reflect"
	"sort"
	"time"

	"bitbucket.org/jazzserve/webapi/utils/pointers/val"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Migration struct {
	ID          *int64
	Name        *string
	Version     int        `gorm:"unique"`
	RunAt       *time.Time `gorm:"default: NOW()"`
	InstalledOn time.Time  `gorm:"default: NOW()"`
	Query       *string    `gorm:"-"`
}

func (m *Migration) Less(anotherM *Migration) bool {
	if m == nil {
		return true
	}
	if anotherM == nil {
		return false
	}
	return m.Version <= anotherM.Version
}

type AutoMigrator interface {
	AutoMigrate(db *gorm.DB) error
}

func (s *storage) Migrate(auto []interface{}, manual []*Migration) error {
	operation := func(tx *gorm.DB) (err error) {
		if err = autoMigrate(tx, auto); err != nil {
			return
		}
		if err = manualMigrate(tx, manual); err != nil {
			return
		}
		return
	}
	return s.RunInTransaction(operation)
}

func autoMigrate(tx *gorm.DB, modelList []interface{}) error {
	for _, t := range modelList {
		tx = tx.AutoMigrate(t)
		if tx.Error != nil {
			tp := reflect.TypeOf(t)
			return errors.Errorf("can't create/modify table %s error: %v", tp.Elem().Name(), tx.Error)
		}
		if mf, ok := t.(AutoMigrator); ok {
			if err := mf.AutoMigrate(tx); err != nil {
				tp := reflect.TypeOf(t)
				return errors.Errorf("%s.AutoMigrate() error: %v", tp.Elem().Name(), err)
			}
		}
	}
	return nil
}

type migrations []*Migration

func (m migrations) Len() int {
	return len(m)
}

func (m migrations) Less(i, j int) bool {
	return m[i].Less(m[j])
}

func (m migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func manualMigrate(tx *gorm.DB, ms []*Migration) error {
	sort.Sort(migrations(ms))

	var lastMigration = &Migration{}

	if err := tx.Last(lastMigration).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return errors.Wrap(err, "can't get the last migration")
		}
		lastMigration = nil
	}

	for _, migration := range ms {
		if migration.Less(lastMigration) {
			continue
		}

		if err := tx.Exec(val.Str(migration.Query)).Error; err != nil {
			return errors.Wrapf(err, "can't exec migration %d_%s", migration.Version, val.Str(migration.Name))
		}

		if err := tx.Create(migration).Error; err != nil {
			return errors.Wrapf(err, "can't create migration %d_%s", migration.Version, val.Str(migration.Name))
		}
	}

	return nil
}

func GetModels() []interface{} {
	return []interface{}{
		&Migration{},
		&Account{},
		&Store{},
		&AccountStore{},
	}
}
