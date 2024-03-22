package grpc

import (
	"context"

	"git.solsynth.dev/hydrogen/identity/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"git.solsynth.dev/hydrogen/identity/pkg/services"
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

	if err := services.NewNotification(client, user, in.Subject, in.Content, links, in.IsImportant); err != nil {
		return nil, err
	}

	return &proto.NotifyReply{IsSent: true}, nil
}
