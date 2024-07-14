package grpc

import (
	"google.golang.org/grpc/reflection"
	"net"

	exproto "git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

import health "google.golang.org/grpc/health/grpc_health_v1"

type Server struct {
	exproto.UnimplementedAuthServer
	proto.UnimplementedNotifyServer
	proto.UnimplementedFriendshipsServer
	proto.UnimplementedRealmsServer
	health.UnimplementedHealthServer

	srv *grpc.Server
}

func NewServer() *Server {
	server := &Server{
		srv: grpc.NewServer(),
	}

	exproto.RegisterAuthServer(server.srv, &Server{})
	proto.RegisterNotifyServer(server.srv, &Server{})
	proto.RegisterFriendshipsServer(server.srv, &Server{})
	proto.RegisterRealmsServer(server.srv, &Server{})
	health.RegisterHealthServer(server.srv, &Server{})

	reflection.Register(server.srv)

	return server
}

func (v *Server) Listen() error {
	listener, err := net.Listen("tcp", viper.GetString("grpc_bind"))
	if err != nil {
		return err
	}

	return v.srv.Serve(listener)
}
