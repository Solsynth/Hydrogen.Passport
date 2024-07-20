package grpc

import (
	"context"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
)

func (v *Server) EmitStreamEvent(ctx context.Context, request *proto.StreamEventRequest) (*proto.StreamEventResponse, error) {
	switch request.GetEvent() {
	case "ClientRegister":
		// No longer need update user online status
		// Based on realtime sever connection status
		break
	case "ClientUnregister":
		// Update user last seen at
		_ = services.SetAccountLastSeen(uint(request.GetUserId()))
	}

	return &proto.StreamEventResponse{}, nil
}
