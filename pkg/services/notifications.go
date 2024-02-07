package services

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/external"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"context"
	"firebase.google.com/go/messaging"
	"github.com/rs/zerolog/log"
)

func AddNotifySubscriber(user models.Account, provider, device, ua string) (models.NotificationSubscriber, error) {
	subscriber := models.NotificationSubscriber{
		UserAgent: ua,
		Provider:  provider,
		DeviceID:  device,
		AccountID: user.ID,
	}

	err := database.C.Save(&subscriber).Error

	return subscriber, err
}

func NewNotification(
	user models.ThirdClient,
	target models.Account,
	subject, content string,
	links []models.NotificationLink,
	important bool,
) error {
	notification := models.Notification{
		Subject:     subject,
		Content:     content,
		Links:       links,
		IsImportant: important,
		ReadAt:      nil,
		SenderID:    &user.ID,
		RecipientID: target.ID,
	}

	if err := database.C.Save(&notification).Error; err != nil {
		return err
	}

	var subscribers []models.NotificationSubscriber
	if err := database.C.Where(&models.NotificationSubscriber{
		AccountID: user.ID,
	}).Find(&subscribers).Error; err != nil {
		// I don't know why cannot get subscribers list, but whatever, the notifications has created
		log.Error().Err(err).Msg("Unexpected error occurred during the notification.")
		return nil
	}

	for _, subscriber := range subscribers {
		switch subscriber.Provider {
		case models.NotifySubscriberFirebase:
			if external.Fire == nil {
				// Didn't configure for firebase support
				break
			}

			ctx := context.Background()
			client, err := external.Fire.Messaging(ctx)
			if err != nil {
				log.Warn().Err(err).Msg("An error occurred when getting firebase FCM client...")
				break
			}

			message := &messaging.Message{
				Notification: &messaging.Notification{
					Title: notification.Subject,
					Body:  notification.Content,
				},
				Token: subscriber.DeviceID,
			}

			if _, err = client.Send(ctx, message); err != nil {
				log.Warn().Err(err).Msg("An error occurred when notify subscriber though firebase FCM...")
			}
		}
	}

	return nil
}
