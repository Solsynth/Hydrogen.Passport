package grpc

import (
	"context"

	exproto "git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	jsoniter "github.com/json-iterator/go"
)

func (v *Server) Authenticate(_ context.Context, in *exproto.AuthRequest) (*exproto.AuthReply, error) {
	ctx, perms, atk, rtk, err := services.Authenticate(in.GetAccessToken(), in.GetRefreshToken(), 0)
	if err != nil {
		return &exproto.AuthReply{
			IsValid: false,
		}, nil
	} else {
		user := ctx.Account
		rawPerms, _ := jsoniter.Marshal(perms)

		userinfo := &exproto.UserInfo{
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

		return &exproto.AuthReply{
			IsValid: true,
			Info: &exproto.AuthInfo{
				NewAccessToken:  &atk,
				NewRefreshToken: &rtk,
				Permissions:     rawPerms,
				TicketId:        uint64(ctx.Ticket.ID),
				Info:            userinfo,
			},
		}, nil
	}
}

func (v *Server) EnsurePermGranted(_ context.Context, in *exproto.CheckPermRequest) (*exproto.CheckPermReply, error) {
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

	return &exproto.CheckPermReply{
		IsValid: valid,
	}, nil
}
