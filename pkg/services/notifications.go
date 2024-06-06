package services

import (
	"context"
	"firebase.google.com/go/messaging"
	"reflect"

	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"github.com/rs/zerolog/log"
)

func AddNotifySubscriber(user models.Account, provider, id, tk, ua string) (models.NotificationSubscriber, error) {
	var prev models.NotificationSubscriber
	var subscriber models.NotificationSubscriber
	if err := database.C.Where(&models.NotificationSubscriber{
		DeviceID:  id,
		AccountID: user.ID,
	}); err != nil {
		subscriber = models.NotificationSubscriber{
			UserAgent:   ua,
			Provider:    provider,
			DeviceID:    id,
			DeviceToken: tk,
			AccountID:   user.ID,
		}
	} else {
		prev = subscriber
	}

	subscriber.UserAgent = ua
	subscriber.Provider = provider
	subscriber.DeviceToken = tk

	var err error
	if !reflect.DeepEqual(subscriber, prev) {
		err = database.C.Save(&subscriber).Error
	}

	return subscriber, err
}

func NewNotification(notification models.Notification) error {
	if err := database.C.Save(&notification).Error; err != nil {
		return err
	}

	go func() {
		err := PushNotification(notification)
		if err != nil {
			log.Error().Err(err).Msg("Unexpected error occurred during the notification.")
		}
	}()

	return nil
}

func PushNotification(notification models.Notification) error {
	frontendAvailable := false
	for conn := range wsConn[notification.RecipientID] {
		frontendAvailable = true
		_ = conn.WriteMessage(1, models.UnifiedCommand{
			Action:  "notifications.new",
			Payload: notification,
		}.Marshal())
	}

	// Skip push notification when frontend notify is available
	if frontendAvailable {
		return nil
	}

	var subscribers []models.NotificationSubscriber
	if err := database.C.Where(&models.NotificationSubscriber{
		AccountID: notification.RecipientID,
	}).Find(&subscribers).Error; err != nil {
		return err
	}

	for _, subscriber := range subscribers {
		switch subscriber.Provider {
		case models.NotifySubscriberFirebase:
			if ExtFire != nil {
				ctx := context.Background()
				client, err := ExtFire.Messaging(ctx)
				if err != nil {
					log.Warn().Err(err).Msg("An error occurred when creating FCM client...")
					break
				}

				message := &messaging.Message{
					Notification: &messaging.Notification{
						Title: notification.Subject,
						Body:  notification.Content,
					},
					Token: subscriber.DeviceToken,
				}

				if response, err := client.Send(ctx, message); err != nil {
					log.Warn().Err(err).Msg("An error occurred when notify subscriber via FCM...")
				} else {
					log.Debug().
						Str("response", response).
						Int("subscriber", int(subscriber.ID)).
						Msg("Notified subscriber via FCM.")
				}
			}
		}
	}

	return nil
}
