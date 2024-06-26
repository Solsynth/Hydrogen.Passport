package grpc

import (
	"context"
	jsoniter "github.com/json-iterator/go"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hydrogen/passport/pkg/proto"
	"github.com/samber/lo"
)

func (v *Server) NotifyUser(_ context.Context, in *proto.NotifyRequest) (*proto.NotifyReply, error) {
	client, err := services.GetThirdClientWithSecret(in.GetClientId(), in.GetClientSecret())
	if err != nil {
		return nil, err
	}

	var user models.Account
	if user, err = services.GetAccount(uint(in.GetRecipientId())); err != nil {
		return nil, err
	}

	var metadata map[string]any
	_ = jsoniter.Unmarshal(in.GetMetadata(), &metadata)

	links := lo.Map(in.GetLinks(), func(item *proto.NotifyLink, index int) models.NotificationLink {
		return models.NotificationLink{
			Label: item.Label,
			Url:   item.Url,
		}
	})

	notification := models.Notification{
		Type:        lo.Ternary(len(in.GetType()) > 0, in.GetType(), "common"),
		Subject:     in.GetSubject(),
		Content:     in.GetContent(),
		Metadata:    metadata,
		Links:       links,
		IsRealtime:  in.GetIsRealtime(),
		IsForcePush: in.GetIsForcePush(),
		RecipientID: user.ID,
		SenderID:    &client.ID,
	}

	if in.GetIsRealtime() {
		if err := services.PushNotification(notification); err != nil {
			return nil, err
		}
	} else {
		if err := services.NewNotification(notification); err != nil {
			return nil, err
		}
	}

	return &proto.NotifyReply{IsSent: true}, nil
}
