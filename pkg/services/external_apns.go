package services

import (
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
	"github.com/spf13/viper"
)

// ExtAPNS is Apple Notification Services client
var ExtAPNS *apns2.Client

func SetupAPNS() error {
	authKey, err := token.AuthKeyFromFile(viper.GetString("apns_credentials"))
	if err != nil {
		return err
	}

	ExtAPNS = apns2.NewTokenClient(&token.Token{
		AuthKey: authKey,
		KeyID:   viper.GetString("apns_credentials_key"),
		TeamID:  viper.GetString("apns_credentials_team"),
	})

	return nil
}
