package grpc

import (
	"context"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
)

func (v *Server) RecordEvent(ctx context.Context, request *proto.RecordEventRequest) (*proto.RecordEventResponse, error) {
	var user models.Account
	var err error
	if user, err = services.GetAccount(uint(request.GetUserId())); err != nil {
		return nil, err
	}

	services.AddEvent(user, request.GetAction(), request.GetTarget(), request.GetIp(), request.GetUserAgent())

	return &proto.RecordEventResponse{IsSuccess: true}, nil
}
