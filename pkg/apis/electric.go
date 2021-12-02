package apis

import (
	"context"
	"log"
	"math"
	"math/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
)

func ElectricApi() server.GrpcApi {
	devices := electric.NewRouter()
	settings := electric.NewMemorySettingsRouter()
	devices.Factory = func(name string) (traits.ElectricApiClient, error) {
		log.Printf("Creating ElectricClient(%v)", name)
		device := electric.NewMemoryDevice()
		// seed with a random load
		_, _ = device.UpdateDemand(context.Background(), &electric.UpdateDemandRequest{
			Demand: &traits.ElectricDemand{
				Rating:  60,
				Voltage: 240,
				Current: float32(math.Round(rand.Float64()*40*100) / 100),
			},
		})
		settings.Add(name, electric.WrapMemorySettings(device))
		return electric.Wrap(device), nil
	}
	return server.Collection(devices, settings)
}
