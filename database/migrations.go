package database

import (
	"github.com/lando-ke/card-api/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Deck{}, &models.Card{})
	if err != nil {
		return err
	}

	return nil
}
