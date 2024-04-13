package grpc

import (
	"context"

	"git.solsynth.dev/hydrogen/passport/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
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

	links := lo.Map(in.GetLinks(), func(item *proto.NotifyLink, index int) models.NotificationLink {
		return models.NotificationLink{
			Label: item.Label,
			Url:   item.Url,
		}
	})

	notification := models.Notification{
		Subject:     in.GetSubject(),
		Content:     in.GetContent(),
		Links:       links,
		IsImportant: in.GetIsImportant(),
		IsRealtime:  in.GetIsRealtime(),
		ReadAt:      nil,
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
