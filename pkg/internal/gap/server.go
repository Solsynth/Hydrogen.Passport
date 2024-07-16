package gap

import (
	"fmt"
	"git.solsynth.dev/hydrogen/dealer/pkg/hyper"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/spf13/viper"
)

var H *hyper.HyperConn

func RegisterService() error {
	grpcBind := strings.SplitN(viper.GetString("grpc_bind"), ":", 2)
	httpBind := strings.SplitN(viper.GetString("bind"), ":", 2)

	outboundIp, _ := GetOutboundIP()

	grpcOutbound := fmt.Sprintf("%s:%s", outboundIp, grpcBind[1])
	httpOutbound := fmt.Sprintf("%s:%s", outboundIp, httpBind[1])

	var err error
	H, err = hyper.NewHyperConn(viper.GetString("dealer.addr"), &proto.ServiceInfo{
		Id:       viper.GetString("id"),
		Type:     hyper.ServiceTypeAuthProvider,
		Label:    "Passport",
		GrpcAddr: grpcOutbound,
		HttpAddr: &httpOutbound,
	})
	if err == nil {
		go func() {
			err := H.KeepRegisterService()
			if err != nil {
				log.Error().Err(err).Msg("An error occurred while registering service...")
			}
		}()
	}

	return err
}
