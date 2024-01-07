package services

import (
	"fmt"

	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
)

func GetAccount(id uint) (models.Account, error) {
	var account models.Account
	if err := database.C.Where(models.Account{
		BaseModel: models.BaseModel{ID: id},
	}).First(&account).Error; err != nil {
		return account, err
	}

	return account, nil
}

func LookupAccount(id string) (models.Account, error) {
	var account models.Account
	if err := database.C.Where(models.Account{Name: id}).First(&account).Error; err == nil {
		return account, nil
	}

	var contact models.AccountContact
	if err := database.C.Where(models.AccountContact{Content: id}).First(&contact).Error; err == nil {
		if err := database.C.
			Where(models.Account{
				BaseModel: models.BaseModel{ID: contact.AccountID},
			}).First(&contact).Error; err == nil {
			return account, err
		}
	}

	return account, fmt.Errorf("account was not found")
}
