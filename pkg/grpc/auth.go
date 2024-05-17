package grpc

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"

	"git.solsynth.dev/hydrogen/passport/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"github.com/spf13/viper"
)

func (v *Server) Authenticate(_ context.Context, in *proto.AuthRequest) (*proto.AuthReply, error) {
	user, perms, atk, rtk, err := services.Authenticate(in.GetAccessToken(), in.GetRefreshToken(), 0)
	if err != nil {
		return &proto.AuthReply{
			IsValid: false,
		}, nil
	} else {
		rawPerms, _ := jsoniter.Marshal(perms)
		return &proto.AuthReply{
			IsValid:      true,
			AccessToken:  &atk,
			RefreshToken: &rtk,
			Permissions:  rawPerms,
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
