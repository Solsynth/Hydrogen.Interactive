package grpc

import (
	pwpb "git.solsynth.dev/hydrogen/identity/pkg/grpc/proto"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var Notify pwpb.NotifyClient
var Auth pwpb.AuthClient

func ConnectPassport() error {
	addr := viper.GetString("identity.grpc_endpoint")
	if conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return err
	} else {
		Notify = pwpb.NewNotifyClient(conn)
		Auth = pwpb.NewAuthClient(conn)
	}

	return nil
}
