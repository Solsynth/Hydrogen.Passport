package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/samber/lo"
)

func GetUserAccountGroup(user models.Account) ([]models.AccountGroup, error) {
	var members []models.AccountGroupMember
	if err := database.C.Where(&models.AccountGroupMember{
		AccountID: user.ID,
	}).Find(&members).Error; err != nil {
		return nil, err
	}

	var groups []models.AccountGroup
	if err := database.C.Where("id IN ?", lo.Map(groups, func(item models.AccountGroup, index int) uint {
		return item.ID
	})).Find(&groups).Error; err != nil {
		return nil, err
	}

	return groups, nil
}
