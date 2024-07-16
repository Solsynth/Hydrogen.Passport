package grpc

import (
	"google.golang.org/grpc/reflection"
	"net"

	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

import health "google.golang.org/grpc/health/grpc_health_v1"

type Server struct {
	proto.UnimplementedAuthServer
	proto.UnimplementedNotifierServer
	proto.UnimplementedRealmServer
	health.UnimplementedHealthServer

	srv *grpc.Server
}

func NewServer() *Server {
	server := &Server{
		srv: grpc.NewServer(),
	}

	proto.RegisterAuthServer(server.srv, server)
	proto.RegisterNotifierServer(server.srv, server)
	proto.RegisterRealmServer(server.srv, server)
	health.RegisterHealthServer(server.srv, server)

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
