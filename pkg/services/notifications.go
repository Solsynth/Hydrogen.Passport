package services

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
)

func AddNotifySubscriber(user models.Account, provider, device, ua string) (models.NotificationSubscriber, error) {
	subscriber := models.NotificationSubscriber{
		UserAgent: ua,
		Provider:  provider,
		DeviceID:  ua,
		AccountID: user.ID,
	}

	err := database.C.Save(&subscriber).Error

	return subscriber, err
}

func NewNotification(user models.ThirdClient, target models.Account, subject, content string, important bool) error {
	notification := models.Notification{
		Subject:     subject,
		Content:     content,
		IsImportant: important,
		ReadAt:      nil,
		SenderID:    &user.ID,
		RecipientID: target.ID,
	}

	if err := database.C.Save(&notification).Error; err != nil {
		return err
	}

	// TODO Notify all the listeners

	return nil
}
