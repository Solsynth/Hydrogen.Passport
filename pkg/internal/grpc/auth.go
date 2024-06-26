package grpc

import (
	"context"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hydrogen/passport/pkg/proto"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
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

		userinfo := &proto.Userinfo{
			Id:          uint64(user.ID),
			Name:        user.Name,
			Nick:        user.Nick,
			Email:       user.GetPrimaryEmail().Content,
			Description: &user.Description,
		}

		if user.Avatar != nil {
			userinfo.Avatar = *user.GetAvatar()
		}
		if user.Banner != nil {
			userinfo.Banner = *user.GetBanner()
		}

		return &proto.AuthReply{
			IsValid:      true,
			AccessToken:  &atk,
			RefreshToken: &rtk,
			Permissions:  rawPerms,
			TicketId:     lo.ToPtr(uint64(ctx.Ticket.ID)),
			Userinfo:     userinfo,
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
