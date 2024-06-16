package gap

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

func Register() error {
	cfg := api.DefaultConfig()
	cfg.Address = viper.GetString("consul.addr")

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	bind := strings.SplitN(viper.GetString("consul.srv_serve"), ":", 2)
	baseAddr := viper.GetString("consul.srv_http")

	port, _ := strconv.Atoi(bind[1])

	registration := new(api.AgentServiceRegistration)
	registration.ID = viper.GetString("id")
	registration.Name = "Hydrogen.Passport"
	registration.Address = bind[0]
	registration.Port = port
	registration.Check = &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("%s/.well-known", baseAddr),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	return client.Agent().ServiceRegister(registration)
}
