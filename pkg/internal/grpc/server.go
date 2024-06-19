package grpc

import (
	"google.golang.org/grpc/reflection"
	"net"

	"git.solsynth.dev/hydrogen/passport/pkg/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

import health "google.golang.org/grpc/health/grpc_health_v1"

type Server struct {
	proto.UnimplementedAuthServer
	proto.UnimplementedNotifyServer
	proto.UnimplementedFriendshipsServer
	proto.UnimplementedRealmsServer
	health.UnimplementedHealthServer
}

var S *grpc.Server

func NewGRPC() {
	S = grpc.NewServer()

	proto.RegisterAuthServer(S, &Server{})
	proto.RegisterNotifyServer(S, &Server{})
	proto.RegisterFriendshipsServer(S, &Server{})
	proto.RegisterRealmsServer(S, &Server{})
	health.RegisterHealthServer(S, &Server{})

	reflection.Register(S)
}

func ListenGRPC() error {
	listener, err := net.Listen("tcp", viper.GetString("grpc_bind"))
	if err != nil {
		return err
	}

	return S.Serve(listener)
}
