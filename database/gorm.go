package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("card-api.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
