package hyper

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"time"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

type HyperConn struct {
	Addr string

	cacheGrpcConn map[string]*grpc.ClientConn
}

func NewHyperConn(addr string) *HyperConn {
	return &HyperConn{
		Addr: addr,

		cacheGrpcConn: make(map[string]*grpc.ClientConn),
	}
}

func (v *HyperConn) DiscoverServiceGRPC(name string) (*grpc.ClientConn, error) {
	if val, ok := v.cacheGrpcConn[name]; ok {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if _, err := health.NewHealthClient(val).Check(ctx, &health.HealthCheckRequest{
			Service: name,
		}); err == nil {
			return val, nil
		} else {
			delete(v.cacheGrpcConn, name)
		}
	}

	target := fmt.Sprintf("consul://%s/%s", v.Addr, name)
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err == nil {
		v.cacheGrpcConn[name] = conn
	}
	return conn, err
}
