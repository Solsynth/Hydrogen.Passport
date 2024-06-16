package grpc

import (
	pcpb "git.solsynth.dev/hydrogen/paperclip/pkg/grpc/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Attachments pcpb.AttachmentsClient

func ConnectPaperclip() error {
	addr := viper.GetString("paperclip.grpc_endpoint")
	if conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return err
	} else {
		Attachments = pcpb.NewAttachmentsClient(conn)
	}

	return nil
}
