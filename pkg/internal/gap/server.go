package gap

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"strings"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

var C *api.Client

func Register() error {
	cfg := api.DefaultConfig()
	cfg.Address = viper.GetString("consul.addr")

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	httpBind := strings.SplitN(viper.GetString("bind"), ":", 2)
	grpcBind := strings.SplitN(viper.GetString("grpc_bind"), ":", 2)

	outboundIp, _ := GetOutboundIP()
	port, _ := strconv.Atoi(httpBind[1])

	registration := new(api.AgentServiceRegistration)
	registration.ID = viper.GetString("id")
	registration.Name = "Hydrogen.Passport"
	registration.Address = outboundIp.String()
	registration.Port = port
	registration.Check = &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%s", outboundIp, grpcBind[1]),
		Timeout:                        "5s",
		Interval:                       "1m",
		DeregisterCriticalServiceAfter: "3m",
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		return err
	} else {
		C = client
		return nil
	}
}

func DiscoverPaperclip() (*grpc.ClientConn, error) {
	return grpc.NewClient(
		"consul://127.0.0.1:8500/Hydrogen.Paperclip",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
}
