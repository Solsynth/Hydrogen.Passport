package database

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"gorm.io/gorm"
)

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(
		&models.Account{},
		&models.AuthFactor{},
		&models.AccountContact{},
		&models.AuthSession{},
		&models.AuthChallenge{},
	); err != nil {
		return err
	}

	return nil
}
