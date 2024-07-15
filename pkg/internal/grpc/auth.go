package grpc

import (
	"context"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"

	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	jsoniter "github.com/json-iterator/go"
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

		userinfo := &proto.UserInfo{
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
			IsValid: true,
			Info: &proto.AuthInfo{
				NewAccessToken:  &atk,
				NewRefreshToken: &rtk,
				Permissions:     rawPerms,
				TicketId:        uint64(ctx.Ticket.ID),
				Info:            userinfo,
			},
		}, nil
	}
}

func (v *Server) EnsurePermGranted(_ context.Context, in *proto.CheckPermRequest) (*proto.CheckPermResponse, error) {
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

	return &proto.CheckPermResponse{
		IsValid: valid,
	}, nil
}

func (v *Server) EnsureUserPermGranted(_ context.Context, in *proto.CheckUserPermRequest) (*proto.CheckUserPermResponse, error) {
	relation, err := services.GetRelationWithTwoNode(uint(in.GetUserId()), uint(in.GetOtherId()))
	if err != nil {
		return &proto.CheckUserPermResponse{
			IsValid: false,
		}, nil
	}

	defaultPerm := relation.Status == models.RelationshipFriend

	var value any
	_ = jsoniter.Unmarshal(in.GetValue(), &value)
	valid := services.HasPermNodeWithDefault(relation.PermNodes, in.GetKey(), value, defaultPerm)

	return &proto.CheckUserPermResponse{
		IsValid: valid,
	}, nil
}
