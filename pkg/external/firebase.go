package external

import (
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var Fire *firebase.App

func SetupFirebase(credentials string) error {
	opt := option.WithCredentialsFile(credentials)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	} else {
		Fire = app
	}

	return nil
}
