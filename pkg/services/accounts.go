package services

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"code.smartsheep.studio/hydrogen/passport/pkg/security"
	"fmt"
	"github.com/samber/lo"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
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
			}).First(&account).Error; err == nil {
			return account, err
		}
	}

	return account, fmt.Errorf("account was not found")
}

func CreateAccount(name, nick, email, password string) (models.Account, error) {
	user := models.Account{
		Name:  name,
		Nick:  nick,
		State: models.PendingAccountState,
		Profile: models.AccountProfile{
			Experience: 100,
		},
		Factors: []models.AuthFactor{
			{
				Type:   models.PasswordAuthFactor,
				Secret: security.HashPassword(password),
			},
		},
		Contacts: []models.AccountContact{
			{
				Type:       models.EmailAccountContact,
				Content:    email,
				VerifiedAt: nil,
			},
		},
		Permissions: datatypes.NewJSONType(make([]string, 0)),
		ConfirmedAt: nil,
	}

	if err := database.C.Create(&user).Error; err != nil {
		return user, err
	}

	if tk, err := NewMagicToken(models.ConfirmMagicToken, &user, nil); err != nil {
		return user, err
	} else if err := NotifyMagicToken(tk); err != nil {
		return user, err
	}

	return user, nil
}

func ConfirmAccount(code string) error {
	var token models.MagicToken
	if err := database.C.Where(&models.MagicToken{
		Code: code,
		Type: models.ConfirmMagicToken,
	}).First(&token).Error; err != nil {
		return err
	} else if token.AssignTo == nil {
		return fmt.Errorf("account was not found")
	}

	var user models.Account
	if err := database.C.Where(&models.Account{
		BaseModel: models.BaseModel{ID: *token.AssignTo},
	}).First(&user).Error; err != nil {
		return err
	}

	return database.C.Transaction(func(tx *gorm.DB) error {
		user.ConfirmedAt = lo.ToPtr(time.Now())

		if err := database.C.Delete(&token).Error; err != nil {
			return err
		}
		if err := database.C.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})
}
