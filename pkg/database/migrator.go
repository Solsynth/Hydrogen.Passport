package database

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"gorm.io/gorm"
)

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(
		&models.Account{},
		&models.AuthFactor{},
		&models.AccountProfile{},
		&models.AccountContact{},
		&models.AuthSession{},
		&models.AuthChallenge{},
		&models.MagicToken{},
		&models.ThirdClient{},
		&models.ActionEvent{},
		&models.Notification{},
	); err != nil {
		return err
	}

	return nil
}
