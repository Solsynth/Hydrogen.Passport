package database

import (
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"gorm.io/gorm"
)

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(
		&models.Account{},
		&models.AuthFactor{},
		&models.AccountProfile{},
		&models.AccountPage{},
		&models.AccountContact{},
		&models.AuthSession{},
		&models.AuthChallenge{},
		&models.MagicToken{},
		&models.ThirdClient{},
		&models.ActionEvent{},
		&models.Notification{},
		&models.NotificationSubscriber{},
	); err != nil {
		return err
	}

	return nil
}
