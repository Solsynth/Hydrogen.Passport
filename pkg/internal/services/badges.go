package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
)

func GrantBadge(user models.Account, badge models.Badge) error {
	badge.AccountID = user.ID
	return database.C.Save(badge).Error
}

func RevokeBadge(badge models.Badge) error {
	return database.C.Delete(&badge).Error
}
