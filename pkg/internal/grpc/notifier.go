package grpc

import (
	"context"
	"fmt"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"

	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
)

func (v *Server) NotifyUser(_ context.Context, in *proto.NotifyUserRequest) (*proto.NotifyResponse, error) {
	var err error
	var user models.Account
	if user, err = services.GetAccount(uint(in.GetUserId())); err != nil {
		return nil, fmt.Errorf("unable to get account: %v", err)
	}

	var metadata map[string]any
	_ = jsoniter.Unmarshal(in.GetNotify().GetMetadata(), &metadata)

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

	var metadata map[string]any
	_ = jsoniter.Unmarshal(in.GetNotify().GetMetadata(), &metadata)

	var notifications []models.Notification
	for _, user := range users {
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

		notifications = append(notifications, notification)
	}

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

	var metadata map[string]any
	_ = jsoniter.Unmarshal(in.GetMetadata(), &metadata)

	var notifications []models.Notification
	for _, user := range users {
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

		notifications = append(notifications, notification)
	}

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
