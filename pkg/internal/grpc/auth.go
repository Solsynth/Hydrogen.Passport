package grpc

import (
	"context"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/samber/lo"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	jsoniter "github.com/json-iterator/go"

	"git.solsynth.dev/hypernet/nexus/pkg/proto"
)

type authenticateServer struct {
	proto.UnimplementedAuthServiceServer
}

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

		if user.AffiliatedID != nil {
			userinfo.AffiliatedTo = lo.ToPtr(uint64(*user.AffiliatedID))
		}
		if user.AutomatedID != nil {
			userinfo.AutomatedBy = lo.ToPtr(uint64(*user.AutomatedID))
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

func (v *Server) ListUserFriends(_ context.Context, in *proto.ListUserRelativeRequest) (*proto.ListUserRelativeResponse, error) {
	tx := database.C.Preload("Account").Where("status = ?", models.RelationshipFriend)

	if in.GetIsRelated() {
		tx = tx.Where("related_id = ?", in.GetUserId())
	} else {
		tx = tx.Where("account_id = ?", in.GetUserId())
	}

	var data []models.AccountRelationship
	if err := tx.Find(&data).Error; err != nil {
		return nil, err
	}

	return &proto.ListUserRelativeResponse{
		Data: lo.Map(data, func(item models.AccountRelationship, index int) *proto.SimpleUserInfo {
			val := &proto.SimpleUserInfo{
				Id:   uint64(item.AccountID),
				Name: item.Account.Name,
				Nick: item.Account.Nick,
			}

			if item.Account.AffiliatedID != nil {
				val.AffiliatedTo = lo.ToPtr(uint64(*item.Account.AffiliatedID))
			}
			if item.Account.AutomatedID != nil {
				val.AutomatedBy = lo.ToPtr(uint64(*item.Account.AutomatedID))
			}

			return val
		}),
	}, nil
}

func (v *Server) ListUserBlocklist(_ context.Context, in *proto.ListUserRelativeRequest) (*proto.ListUserRelativeResponse, error) {
	tx := database.C.Preload("Account").Where("status = ?", models.RelationshipBlocked)

	if in.GetIsRelated() {
		tx = tx.Where("related_id = ?", in.GetUserId())
	} else {
		tx = tx.Where("account_id = ?", in.GetUserId())
	}

	var data []models.AccountRelationship
	if err := tx.Find(&data).Error; err != nil {
		return nil, err
	}

	return &proto.ListUserRelativeResponse{
		Data: lo.Map(data, func(item models.AccountRelationship, index int) *proto.SimpleUserInfo {
			val := &proto.SimpleUserInfo{
				Id:   uint64(item.AccountID),
				Name: item.Account.Name,
				Nick: item.Account.Nick,
			}

			if item.Account.AffiliatedID != nil {
				val.AffiliatedTo = lo.ToPtr(uint64(*item.Account.AffiliatedID))
			}
			if item.Account.AutomatedID != nil {
				val.AutomatedBy = lo.ToPtr(uint64(*item.Account.AutomatedID))
			}

			return val
		}),
	}, nil
}
