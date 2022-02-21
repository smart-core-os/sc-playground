package powersupply

import (
	"math"
	"math/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/powersupply"
	"github.com/smart-core-os/sc-playground/pkg/node"
)

func Activate(n *node.Node) {
	settings := powersupply.NewMemorySettingsApiRouter()
	devices := powersupply.NewApiRouter(
		powersupply.WithPowerSupplyApiClientFactory(func(name string) (traits.PowerSupplyApiClient, error) {
			device := powersupply.NewMemoryDevice()
			// seed with a random load
			device.SetLoad(float32(math.Round(rand.Float64()*40*100) / 100))
			settings.Add(name, powersupply.WrapMemorySettingsApi(device))
			n.Announce(name, node.HasTrait(trait.PowerSupply))
			return powersupply.WrapApi(device), nil
		}),
	)
	n.AddRouter(devices, settings)
}
