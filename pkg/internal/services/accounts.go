package services

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/datatypes"
	"time"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
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

func GetAccountWithName(alias string) (models.Account, error) {
	var account models.Account
	if err := database.C.Where(models.Account{
		Name: alias,
	}).First(&account).Error; err != nil {
		return account, err
	}

	return account, nil
}

func LookupAccount(probe string) (models.Account, error) {
	var account models.Account
	if err := database.C.Where(models.Account{Name: probe}).First(&account).Error; err == nil {
		return account, nil
	}

	var contact models.AccountContact
	if err := database.C.Where(models.AccountContact{Content: probe}).First(&contact).Error; err == nil {
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
		Name: name,
		Nick: nick,
		Profile: models.AccountProfile{
			Experience: 100,
		},
		Factors: []models.AuthFactor{
			{
				Type:   models.PasswordAuthFactor,
				Secret: HashPassword(password),
			},
			{
				Type:   models.EmailPasswordFactor,
				Secret: uuid.NewString()[:8],
			},
		},
		Contacts: []models.AccountContact{
			{
				Type:       models.EmailAccountContact,
				Content:    email,
				IsPrimary:  true,
				VerifiedAt: nil,
			},
		},
		PermNodes:   datatypes.JSONMap(viper.GetStringMap("permissions.default")),
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
	token, err := ValidateMagicToken(code, models.ConfirmMagicToken)
	if err != nil {
		return err
	}

	var user models.Account
	if err := database.C.Where(&models.Account{
		BaseModel: models.BaseModel{ID: *token.AccountID},
	}).First(&user).Error; err != nil {
		return err
	}

	return database.C.Transaction(func(tx *gorm.DB) error {
		user.ConfirmedAt = lo.ToPtr(time.Now())

		for k, v := range viper.GetStringMap("permissions.verified") {
			if val, ok := user.PermNodes[k]; !ok {
				user.PermNodes[k] = v
			} else if !ComparePermNode(val, v) {
				user.PermNodes[k] = v
			}
		}

		if err := database.C.Delete(&token).Error; err != nil {
			return err
		}
		if err := database.C.Save(&user).Error; err != nil {
			return err
		}

		InvalidAuthCacheWithUser(user.ID)

		return nil
	})
}

func DeleteAccount(id uint) error {
	tx := database.C.Begin()

	for _, model := range []any{
		&models.Badge{},
		&models.RealmMember{},
		&models.AccountContact{},
		&models.AuthFactor{},
		&models.AuthTicket{},
		&models.MagicToken{},
		&models.ThirdClient{},
		&models.Notification{},
		&models.NotificationSubscriber{},
		&models.AccountFriendship{},
	} {
		if err := tx.Delete(model, "account_id = ?", id).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Delete(&models.Account{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func RecycleUnConfirmAccount() {
	var hitList []models.Account
	if err := database.C.Where("confirmed_at IS NULL").Find(&hitList).Error; err != nil {
		log.Error().Err(err).Msg("An error occurred while recycling accounts...")
		return
	}

	if len(hitList) > 0 {
		log.Info().Int("count", len(hitList)).Msg("Going to recycle those un-confirmed accounts...")
		for _, entry := range hitList {
			if err := DeleteAccount(entry.ID); err != nil {
				log.Error().Err(err).Msg("An error occurred while recycling accounts...")
			}
		}
	}
}

func SetAccountLastSeen(uid uint) error {
	var profile models.AccountProfile
	if err := database.C.Where("account_id = ?", uid).First(&profile).Error; err != nil {
		return err
	}

	profile.LastSeenAt = lo.ToPtr(time.Now())

	return database.C.Save(&profile).Error
}
