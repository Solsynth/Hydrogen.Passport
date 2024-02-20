package grpc

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/grpc/proto"
	"code.smartsheep.studio/hydrogen/identity/pkg/models"
	"code.smartsheep.studio/hydrogen/identity/pkg/services"
	"context"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	proto.UnimplementedAuthServer
	proto.UnimplementedNotifyServer
}

func (v *Server) Authenticate(_ context.Context, in *proto.AuthRequest) (*proto.AuthReply, error) {
	user, atk, rtk, err := services.Authenticate(in.GetAccessToken(), in.GetRefreshToken(), 0)
	if err != nil {
		return &proto.AuthReply{
			IsValid: false,
		}, nil
	} else {
		return &proto.AuthReply{
			IsValid:      true,
			AccessToken:  &atk,
			RefreshToken: &rtk,
			Userinfo: &proto.Userinfo{
				Name:        user.Name,
				Nick:        user.Nick,
				Avatar:      user.Avatar,
				Email:       user.GetPrimaryEmail().Content,
				Description: nil,
			},
		}, nil
	}
}

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

func StartGrpc() error {
	listen, err := net.Listen("tcp", viper.GetString("grpc_bind"))
	if err != nil {
		return err
	}

	server := grpc.NewServer()

	proto.RegisterAuthServer(server, &Server{})
	proto.RegisterNotifyServer(server, &Server{})

	reflection.Register(server)

	return server.Serve(listen)
}
