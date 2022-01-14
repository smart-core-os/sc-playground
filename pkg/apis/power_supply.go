package apis

import (
	"log"
	"math"
	"math/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/powersupply"
)

func PowerSupplyApi(traiter Traiter) server.GrpcApi {
	settings := powersupply.NewMemorySettingsApiRouter()
	devices := powersupply.NewApiRouter(
		powersupply.WithPowerSupplyApiClientFactory(func(name string) (traits.PowerSupplyApiClient, error) {
			log.Printf("Creating PowerSupplyClient(%v)", name)
			traiter.Trait(name, trait.PowerSupply)
			device := powersupply.NewMemoryDevice()
			// seed with a random load
			device.SetLoad(float32(math.Round(rand.Float64()*40*100) / 100))
			settings.Add(name, powersupply.WrapMemorySettingsApi(device))
			return powersupply.WrapApi(device), nil
		}),
	)
	return server.Collection(devices, settings)
}
