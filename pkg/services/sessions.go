package services

import (
	"time"

	"git.solsynth.dev/hydrogen/identity/pkg/database"
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func LookupSessionWithToken(tokenId string) (models.AuthSession, error) {
	var session models.AuthSession
	if err := database.C.
		Where(models.AuthSession{AccessToken: tokenId}).
		Or(models.AuthSession{RefreshToken: tokenId}).
		First(&session).Error; err != nil {
		return session, err
	}

	return session, nil
}

func PerformAutoSignoff() *gorm.DB {
	signoffDuration := time.Duration(viper.GetInt64("security.auto_signoff_duration")) * time.Second
	return database.C.Where("last_grant_at < ?", time.Now().Add(-signoffDuration)).Delete(&models.AuthSession{})
}
