package grpc

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"

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
		IsRealtime:  in.GetNotify().GetIsRealtime(),
		IsForcePush: in.GetNotify().GetIsForcePush(),
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
