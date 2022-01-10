package apis

import (
	"log"
	"math"
	"math/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait/powersupply"
)

func PowerSupplyApi() server.GrpcApi {
	devices := powersupply.NewApiRouter()
	settings := powersupply.NewMemorySettingsApiRouter()
	devices.Factory = func(name string) (traits.PowerSupplyApiClient, error) {
		log.Printf("Creating PowerSupplyClient(%v)", name)
		device := powersupply.NewMemoryDevice()
		// seed with a random load
		device.SetLoad(float32(math.Round(rand.Float64()*40*100) / 100))
		settings.Add(name, powersupply.WrapMemorySettingsApi(device))
		return powersupply.WrapApi(device), nil
	}
	return server.Collection(devices, settings)
}
