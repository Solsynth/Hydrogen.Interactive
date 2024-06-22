package gap

import (
	"git.solsynth.dev/hydrogen/passport/pkg/hyper"
	"github.com/spf13/viper"
)

var H *hyper.HyperConn

func NewHyperClient() {
	H = hyper.NewHyperConn(viper.GetString("consul.addr"))
}
