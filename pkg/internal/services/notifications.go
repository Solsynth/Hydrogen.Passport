package services

import (
	"context"
	"fmt"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/gap"
	"reflect"
	"time"

	"firebase.google.com/go/messaging"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/sideshow/apns2"
	payload2 "github.com/sideshow/apns2/payload"
	"github.com/spf13/viper"
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

// NewNotification will create a notification and push via the push method it
func NewNotification(notification models.Notification) error {
	if err := database.C.Save(&notification).Error; err != nil {
		return err
	}

	if err := PushNotification(notification); err != nil {
		return err
	}

	return nil
}

// PushNotification will push the notification whatever it exists record in the
// database Recommend pushing another goroutine when you need to push a lot of
// notifications And just use a block statement when you just push one
// notification.
// The time of creating a new subprocess is much more than push notification.
func PushNotification(notification models.Notification) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := proto.NewStreamControllerClient(gap.H.GetDealerGrpcConn()).PushStream(ctx, &proto.PushStreamRequest{
		UserId: uint64(notification.AccountID),
		Body: models.UnifiedCommand{
			Action:  "notifications.new",
			Payload: notification,
		}.Marshal(),
	})
	if err != nil {
		return fmt.Errorf("failed to push via websocket: %v", err)
	}

	// Skip push notification
	if GetStatusDisturbable(notification.AccountID) != nil {
		return nil
	}

	var subscribers []models.NotificationSubscriber
	if err := database.C.Where(&models.NotificationSubscriber{
		AccountID: notification.AccountID,
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
						Title: notification.Title,
						Body:  notification.Body,
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
		case models.NotifySubscriberAPNs:
			if ExtAPNS != nil {
				data, err := payload2.
					NewPayload().
					AlertTitle(notification.Title).
					AlertBody(notification.Body).
					Sound("default").
					Category(notification.Topic).
					MarshalJSON()
				if err != nil {
					log.Warn().Err(err).Msg("An error occurred when preparing to notify subscriber via APNs...")
				}
				payload := &apns2.Notification{
					ApnsID:      subscriber.DeviceID,
					DeviceToken: subscriber.DeviceToken,
					Topic:       viper.GetString("apns_topic"),
					Payload:     data,
				}

				if resp, err := ExtAPNS.Push(payload); err != nil {
					log.Warn().Err(err).Msg("An error occurred when notify subscriber via APNs...")
				} else {
					log.Debug().
						Str("reason", resp.Reason).
						Int("status", resp.StatusCode).
						Int("subscriber", int(subscriber.ID)).
						Msg("Notified subscriber via APNs.")
				}
			}
		}
	}

	return nil
}
