package grpc

import (
	"context"
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"github.com/samber/lo"
)

func (v *Server) ListFriendship(_ context.Context, request *proto.FriendshipLookupRequest) (*proto.ListFriendshipResponse, error) {
	account, err := services.GetAccount(uint(request.GetAccountId()))
	if err != nil {
		return nil, err
	}
	friends, err := services.ListFriend(account, models.FriendshipStatus(request.GetStatus()))
	if err != nil {
		return nil, err
	}

	return &proto.ListFriendshipResponse{
		Data: lo.Map(friends, func(item models.AccountFriendship, index int) *proto.FriendshipResponse {
			return &proto.FriendshipResponse{
				AccountId: uint64(item.AccountID),
				RelatedId: uint64(item.RelatedID),
				Status:    uint32(item.Status),
			}
		}),
	}, nil
}

func (v *Server) GetFriendship(ctx context.Context, request *proto.FriendshipTwoSideLookupRequest) (*proto.FriendshipResponse, error) {
	friend, err := services.GetFriendWithTwoSides(uint(request.GetAccountId()), uint(request.GetRelatedId()))
	if err != nil {
		return nil, err
	} else if friend.Status != models.FriendshipStatus(request.GetStatus()) {
		return nil, fmt.Errorf("status mismatch")
	}

	return &proto.FriendshipResponse{
		AccountId: uint64(friend.AccountID),
		RelatedId: uint64(friend.RelatedID),
		Status:    uint32(friend.Status),
	}, nil
}
