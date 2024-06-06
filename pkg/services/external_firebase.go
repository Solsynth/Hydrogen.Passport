package services

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

// ExtFire is the firebase app client
var ExtFire *firebase.App

func SetupFirebase() error {
	opt := option.WithCredentialsFile(viper.GetString("firebase_credentials"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	} else {
		ExtFire = app
	}

	return nil
}
