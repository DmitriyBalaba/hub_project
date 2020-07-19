package postgres

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewDB(c *Config, debugMode bool) (db *gorm.DB, err error) {
	if c == nil {
		panic("cannot create new postgres db with nil config")
	}

	db, err = gorm.Open(Driver, c.composeConnectionString())
	if err != nil {
		return
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)

	db.LogMode(debugMode)

	return
}
