package node

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"google.golang.org/grpc"
)

func ExampleNode() {
	n := New("my-server")

	// add supported apis
	n.AddRouter(
		onoff.NewApiRouter(),
		light.NewApiRouter(),
	)

	// announce devices
	n.Announce(
		"my-device",
		HasClient(""),
		HasTrait(trait.OnOff, WithClients(onOffClient)),
		HasTrait(trait.Light, WithClients(lightClient)),
	)

	// register our node with the grpc server
	n.Register(grpcServer)
}

var (
	grpcServer  *grpc.Server
	onOffClient traits.OnOffApiClient
	lightClient traits.LightApiClient
)
