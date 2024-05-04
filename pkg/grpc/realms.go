package grpc

import (
	"context"
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (v *Server) ListCommunityRealm(ctx context.Context, empty *emptypb.Empty) (*proto.ListRealmResponse, error) {
	realms, err := services.ListCommunityRealm()
	if err != nil {
		return nil, err
	}

	return &proto.ListRealmResponse{
		Data: lo.Map(realms, func(item models.Realm, index int) *proto.RealmResponse {
			return &proto.RealmResponse{
				Alias:       item.Alias,
				Name:        item.Name,
				Description: item.Description,
				IsPublic:    item.IsPublic,
				IsCommunity: item.IsCommunity,
			}
		}),
	}, nil
}

func (v *Server) ListAvailableRealm(ctx context.Context, request *proto.RealmLookupWithUserRequest) (*proto.ListRealmResponse, error) {
	account, err := services.GetAccount(uint(request.GetUserId()))
	if err != nil {
		return nil, fmt.Errorf("unable to find target account: %v", err)
	}
	realms, err := services.ListAvailableRealm(account)
	if err != nil {
		return nil, err
	}

	return &proto.ListRealmResponse{
		Data: lo.Map(realms, func(item models.Realm, index int) *proto.RealmResponse {
			return &proto.RealmResponse{
				Alias:       item.Alias,
				Name:        item.Name,
				Description: item.Description,
				IsPublic:    item.IsPublic,
				IsCommunity: item.IsCommunity,
			}
		}),
	}, nil
}

func (v *Server) ListOwnedRealm(ctx context.Context, request *proto.RealmLookupWithUserRequest) (*proto.ListRealmResponse, error) {
	account, err := services.GetAccount(uint(request.GetUserId()))
	if err != nil {
		return nil, fmt.Errorf("unable to find target account: %v", err)
	}
	realms, err := services.ListOwnedRealm(account)
	if err != nil {
		return nil, err
	}

	return &proto.ListRealmResponse{
		Data: lo.Map(realms, func(item models.Realm, index int) *proto.RealmResponse {
			return &proto.RealmResponse{
				Alias:       item.Alias,
				Name:        item.Name,
				Description: item.Description,
				IsPublic:    item.IsPublic,
				IsCommunity: item.IsCommunity,
			}
		}),
	}, nil
}

func (v *Server) GetRealm(ctx context.Context, request *proto.RealmLookupRequest) (*proto.RealmResponse, error) {
	var realm models.Realm

	tx := database.C.Model(&models.Realm{})
	if request.Id != nil {
		tx = tx.Where("id = ?", request.GetId())
	}
	if request.Alias != nil {
		tx = tx.Where("alias = ?", request.GetAlias())
	}
	if request.IsPublic != nil {
		tx = tx.Where("is_public = ?", request.GetIsPublic())
	}
	if request.IsCommunity != nil {
		tx = tx.Where("is_community = ?", request.GetIsCommunity())
	}

	if err := tx.First(&realm).Error; err != nil {
		return nil, err
	}

	return &proto.RealmResponse{
		Alias:       realm.Alias,
		Name:        realm.Name,
		Description: realm.Description,
		IsPublic:    realm.IsPublic,
		IsCommunity: realm.IsCommunity,
	}, nil
}

func (v *Server) ListRealmMember(ctx context.Context, request *proto.RealmMemberLookupRequest) (*proto.ListRealmMemberResponse, error) {
	var members []models.RealmMember
	tx := database.C.Where("realm_id = ?", request.GetRealmId())
	if request.UserId != nil {
		tx = tx.Where("account_id = ?", request.GetUserId())
	}

	if err := tx.Find(&members).Error; err != nil {
		return nil, err
	}

	return &proto.ListRealmMemberResponse{
		Data: lo.Map(members, func(item models.RealmMember, index int) *proto.RealmMemberResponse {
			return &proto.RealmMemberResponse{
				RealmId:    uint64(item.RealmID),
				UserId:     uint64(item.AccountID),
				PowerLevel: int32(item.PowerLevel),
			}
		}),
	}, nil
}

func (v *Server) GetRealmMember(ctx context.Context, request *proto.RealmMemberLookupRequest) (*proto.RealmMemberResponse, error) {
	var member models.RealmMember
	tx := database.C.Where("realm_id = ?", request.GetRealmId())
	if request.UserId != nil {
		tx = tx.Where("account_id = ?", request.GetUserId())
	}

	if err := tx.First(&member).Error; err != nil {
		return nil, err
	}

	return &proto.RealmMemberResponse{
		RealmId:    uint64(member.RealmID),
		UserId:     uint64(member.AccountID),
		PowerLevel: int32(member.PowerLevel),
	}, nil
}
