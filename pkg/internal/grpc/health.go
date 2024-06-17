package grpc

import (
	"context"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

func (v *Server) Check(ctx context.Context, request *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

func (v *Server) Watch(request *health.HealthCheckRequest, server health.Health_WatchServer) error {
	for {
		if server.Send(&health.HealthCheckResponse{
			Status: health.HealthCheckResponse_SERVING,
		}) != nil {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}
