package grpc

import (
	"context"
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"github.com/rs/zerolog/log"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"github.com/samber/lo"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hydrogen/passport/pkg/proto"
)

func (v *Server) NotifyUser(_ context.Context, in *proto.NotifyUserRequest) (*proto.NotifyResponse, error) {
	var err error
	var user models.Account
	if user, err = services.GetAccount(uint(in.GetUserId())); err != nil {
		return nil, fmt.Errorf("unable to get account: %v", err)
	}

	metadata := nex.DecodeMap(in.GetNotify().GetMetadata())

	notification := models.Notification{
		Topic:       in.GetNotify().GetTopic(),
		Title:       in.GetNotify().GetTitle(),
		Subtitle:    in.GetNotify().Subtitle,
		Body:        in.GetNotify().GetBody(),
		Metadata:    metadata,
		Avatar:      in.GetNotify().Avatar,
		Picture:     in.GetNotify().Picture,
		IsRealtime:  in.GetNotify().GetIsRealtime(),
		IsForcePush: in.GetNotify().GetIsForcePush(),
		Account:     user,
		AccountID:   user.ID,
	}

	log.Debug().Str("topic", notification.Topic).Uint("uid", notification.AccountID).Msg("Notifying user...")

	if notification.IsRealtime {
		if err := services.PushNotification(notification); err != nil {
			return nil, err
		}
	} else {
		if err := services.NewNotification(notification); err != nil {
			return nil, err
		}
	}

	return &proto.NotifyResponse{
		IsSuccess: true,
	}, nil
}

func (v *Server) NotifyUserBatch(_ context.Context, in *proto.NotifyUserBatchRequest) (*proto.NotifyResponse, error) {
	var err error
	var users []models.Account
	if users, err = services.GetAccountList(lo.Map(in.GetUserId(), func(item uint64, index int) uint {
		return uint(item)
	})); err != nil {
		return nil, fmt.Errorf("unable to get account: %v", err)
	}

	metadata := nex.DecodeMap(in.GetNotify().GetMetadata())

	var checklist = make(map[uint]bool, len(users))
	var notifications []models.Notification
	for _, user := range users {
		if _, ok := checklist[user.ID]; ok {
			continue
		}

		notification := models.Notification{
			Topic:       in.GetNotify().GetTopic(),
			Title:       in.GetNotify().GetTitle(),
			Subtitle:    in.GetNotify().Subtitle,
			Body:        in.GetNotify().GetBody(),
			Metadata:    metadata,
			Avatar:      in.GetNotify().Avatar,
			Picture:     in.GetNotify().Picture,
			IsRealtime:  in.GetNotify().GetIsRealtime(),
			IsForcePush: in.GetNotify().GetIsForcePush(),
			Account:     user,
			AccountID:   user.ID,
		}
		checklist[user.ID] = true

		notifications = append(notifications, notification)
	}

	log.Debug().Str("topic", notifications[0].Topic).Any("uid", lo.Keys(checklist)).Msg("Notifying users...")

	if in.GetNotify().GetIsRealtime() {
		services.PushNotificationBatch(notifications)
	} else {
		if err := services.NewNotificationBatch(notifications); err != nil {
			return nil, err
		}
	}

	return &proto.NotifyResponse{
		IsSuccess: true,
	}, nil
}

func (v *Server) NotifyAllUser(_ context.Context, in *proto.NotifyRequest) (*proto.NotifyResponse, error) {
	var users []models.Account
	if err := database.C.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("unable to get account: %v", err)
	}

	metadata := nex.DecodeMap(in.GetMetadata())

	var checklist = make(map[uint]bool, len(users))
	var notifications []models.Notification
	for _, user := range users {
		if checklist[user.ID] {
			continue
		}

		notification := models.Notification{
			Topic:       in.GetTopic(),
			Title:       in.GetTitle(),
			Subtitle:    in.Subtitle,
			Body:        in.GetBody(),
			Metadata:    metadata,
			Avatar:      in.Avatar,
			Picture:     in.Picture,
			IsRealtime:  in.GetIsRealtime(),
			IsForcePush: in.GetIsForcePush(),
			Account:     user,
			AccountID:   user.ID,
		}
		checklist[user.ID] = true

		notifications = append(notifications, notification)
	}

	log.Debug().Str("topic", notifications[0].Topic).Any("uid", lo.Keys(checklist)).Msg("Notifying users...")

	if in.GetIsRealtime() {
		services.PushNotificationBatch(notifications)
	} else {
		if err := services.NewNotificationBatch(notifications); err != nil {
			return nil, err
		}
	}

	return &proto.NotifyResponse{
		IsSuccess: true,
	}, nil
}
