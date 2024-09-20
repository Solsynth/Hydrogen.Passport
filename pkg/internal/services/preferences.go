package services

import (
	"errors"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/samber/lo"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func GetNotificationPreference(account models.Account) (models.PreferenceNotification, error) {
	var notification models.PreferenceNotification
	if err := database.C.Where("account_id = ?", account.ID).First(&notification).Error; err != nil {
		return notification, err
	}
	return notification, nil
}

func UpdateNotificationPreference(account models.Account, config map[string]bool) (models.PreferenceNotification, error) {
	var notification models.PreferenceNotification
	var err error
	if notification, err = GetNotificationPreference(account); err != nil {
		notification = models.PreferenceNotification{
			AccountID: account.ID,
			Config: datatypes.JSONMap(
				lo.MapValues(config, func(v bool, k string) any { return v }),
			),
		}
	} else {
		notification.Config = datatypes.JSONMap(
			lo.MapValues(config, func(v bool, k string) any { return v }),
		)
	}

	err = database.C.Save(&notification).Error
	return notification, err
}

func CheckNotificationNotifiable(account models.Account, topic string) bool {
	var notification models.PreferenceNotification
	if err := database.C.Where("account_id = ?", account.ID).First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
		return false
	}
	if val, ok := notification.Config[topic]; ok {
		if status, ok := val.(bool); ok {
			return status
		}
	}
	return true
}

func CheckNotificationNotifiableBatch(accounts []models.Account, topic string) []bool {
	var notifications []models.PreferenceNotification
	if err := database.C.Where("account_id IN ?", accounts).Find(&notifications).Error; err != nil {
		return lo.Map(accounts, func(item models.Account, index int) bool {
			return false
		})
	}

	var notifiable []bool
	for _, notification := range notifications {
		if val, ok := notification.Config[topic]; ok {
			if status, ok := val.(bool); ok {
				notifiable = append(notifiable, status)
				continue
			}
		}
		notifiable = append(notifiable, true)
	}

	return notifiable
}
