package grpc

import (
	"context"
	"fmt"
	"net"

	"git.solsynth.dev/hydrogen/identity/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"git.solsynth.dev/hydrogen/identity/pkg/services"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
				Id:          uint64(user.ID),
				Name:        user.Name,
				Nick:        user.Nick,
				Email:       user.GetPrimaryEmail().Content,
				Avatar:      fmt.Sprintf("https://%s/api/avatar/%s", viper.GetString("domain"), user.Avatar),
				Banner:      fmt.Sprintf("https://%s/api/avatar/%s", viper.GetString("domain"), user.Banner),
				Description: &user.Description,
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
