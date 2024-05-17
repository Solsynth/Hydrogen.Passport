package grpc

import (
	"context"
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func (v *Server) Authenticate(_ context.Context, in *proto.AuthRequest) (*proto.AuthReply, error) {
	ctx, perms, atk, rtk, err := services.Authenticate(in.GetAccessToken(), in.GetRefreshToken(), 0)
	if err != nil {
		return &proto.AuthReply{
			IsValid: false,
		}, nil
	} else {
		user := ctx.Account
		rawPerms, _ := jsoniter.Marshal(perms)
		return &proto.AuthReply{
			IsValid:      true,
			AccessToken:  &atk,
			RefreshToken: &rtk,
			Permissions:  rawPerms,
			TicketId:     lo.ToPtr(uint64(ctx.Ticket.ID)),
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

func (v *Server) CheckPerm(_ context.Context, in *proto.CheckPermRequest) (*proto.CheckPermReply, error) {
	claims, err := services.DecodeJwt(in.GetToken())
	if err != nil {
		return nil, err
	}
	ctx, err := services.GetAuthContext(claims.ID)
	if err != nil {
		return nil, err
	}

	var heldPerms map[string]any
	rawHeldPerms, _ := jsoniter.Marshal(ctx.Account.PermNodes)
	_ = jsoniter.Unmarshal(rawHeldPerms, &heldPerms)

	var value any
	_ = jsoniter.Unmarshal(in.GetValue(), &value)
	perms := services.FilterPermNodes(heldPerms, ctx.Ticket.Claims)
	valid := services.HasPermNode(perms, in.GetKey(), value)

	return &proto.CheckPermReply{
		IsValid: valid,
	}, nil
}
