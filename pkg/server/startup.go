package server

import (
	"code.smartsheep.studio/hydrogen/bus/pkg/kit/publisher"
	"github.com/spf13/viper"
)

const (
	Hostname  = "hydrogen.passport"
	Namespace = "passport"
)

var C *publisher.PublisherConnection

func InitConnection(addr, id string) error {
	if conn, err := publisher.NewConnection(
		addr,
		id,
		Hostname,
		Namespace,
		viper.Get("credentials"),
	); err != nil {
		return err
	} else {
		C = conn
	}

	return nil
}

func PublishCommands(conn *publisher.PublisherConnection) error {
	for k, v := range Commands {
		if err := conn.PublishCommand(k, v); err != nil {
			return err
		}
	}

	return nil
}
