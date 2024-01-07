package services

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
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
