package hyper

import (
	"context"
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/proto"
	"google.golang.org/grpc"
	"time"
)

func (v *HyperConn) DoAuthenticate(atk, rtk string) (acc *proto.Userinfo, accessTk string, refreshTk string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var in *grpc.ClientConn
	in, err = v.DiscoverServiceGRPC("Hydrogen.Passport")
	if err != nil {
		return
	}

	var reply *proto.AuthReply
	reply, err = proto.NewAuthClient(in).Authenticate(ctx, &proto.AuthRequest{
		AccessToken:  atk,
		RefreshToken: &rtk,
	})
	if err != nil {
		return
	}
	if reply != nil {
		acc = reply.GetUserinfo()
		accessTk = reply.GetAccessToken()
		refreshTk = reply.GetRefreshToken()
		if !reply.IsValid {
			err = fmt.Errorf("invalid authorization context")
			return
		}
	}

	return
}

func (v *HyperConn) DoCheckPerm(atk string, key string, val []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	in, err := v.DiscoverServiceGRPC("Hydrogen.Passport")
	if err != nil {
		return err
	}

	reply, err := proto.NewAuthClient(in).CheckPerm(ctx, &proto.CheckPermRequest{
		Token: atk,
		Key:   key,
		Value: val,
	})
	if err != nil {
		return err
	} else if !reply.GetIsValid() {
		return fmt.Errorf("missing permission: %s", key)
	}

	return nil
}
