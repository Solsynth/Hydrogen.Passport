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
