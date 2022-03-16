package powersupply

import (
	"log"
	"math"
	"math/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/powersupply"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"google.golang.org/protobuf/proto"
)

func Activate(n *node.Node) {
	settings := powersupply.NewMemorySettingsApiRouter()
	devices := powersupply.NewApiRouter(
		powersupply.WithPowerSupplyApiClientFactory(func(name string) (traits.PowerSupplyApiClient, error) {
			device := powersupply.NewMemoryDevice()
			// seed with a random load
			device.SetLoad(float32(math.Round(rand.Float64()*40*100) / 100))
			return powersupply.WrapApi(device), nil
		}),
		router.WithOnCommit(func(name string, client interface{}) {
			device, ok := wrap.UnwrapFully(client).(*powersupply.MemoryDevice)
			if !ok {
				return
			}
			log.Printf("PowerSupplyApiClient(%v) auto-created", name)
			settings.Add(name, powersupply.WrapMemorySettingsApi(device))
			n.Announce(name, node.HasTrait(trait.PowerSupply))
		}),
	)
	n.AddRouter(devices, settings)
	n.AddTraitFactory(trait.PowerSupply, func(name string, _ proto.Message) error {
		_, err := devices.Get(name)
		return err
	})
}
