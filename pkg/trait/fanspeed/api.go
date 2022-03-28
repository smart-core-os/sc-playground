package fanspeed

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/fanspeed"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func Activate(n *node.Node) {
	r := fanspeed.NewApiRouter(
		fanspeed.WithFanSpeedApiClientFactory(func(name string) (traits.FanSpeedApiClient, error) {
			return fanspeed.WrapApi(fanspeed.NewModelServer(fanspeed.NewModel())), nil
		}),
		n.AnnounceOnRouterChange(trait.FanSpeed),
	)
	n.AddRouter(r)
	n.AddTraitFactory(trait.FanSpeed, func(name string, _ proto.Message) error {
		_, err := r.Get(name)
		return err
	})
	n.AddClientFactory(trait.FanSpeed, func(conn *grpc.ClientConn) interface{} {
		return traits.NewFanSpeedApiClient(conn)
	})
}
