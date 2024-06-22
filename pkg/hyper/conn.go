package hyper

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HyperConn struct {
	Addr string
}

func NewHyperConn(addr string) *HyperConn {
	return &HyperConn{Addr: addr}
}

func (v *HyperConn) DiscoverServiceGRPC(name string) (*grpc.ClientConn, error) {
	target := fmt.Sprintf("consul://%s/%s", v.Addr, name)
	return grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
}
