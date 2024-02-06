package services

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
)

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
