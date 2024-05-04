package grpc

import (
	idpb "git.solsynth.dev/hydrogen/passport/pkg/grpc/proto"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var Realms idpb.RealmsClient
var Friendships idpb.FriendshipsClient
var Notify idpb.NotifyClient
var Auth idpb.AuthClient

func ConnectPassport() error {
	addr := viper.GetString("passport.grpc_endpoint")
	if conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return err
	} else {
		Realms = idpb.NewRealmsClient(conn)
		Friendships = idpb.NewFriendshipsClient(conn)
		Notify = idpb.NewNotifyClient(conn)
		Auth = idpb.NewAuthClient(conn)
	}

	return nil
}
