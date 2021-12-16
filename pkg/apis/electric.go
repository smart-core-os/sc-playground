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
		var voltage float32 = 240
		var rating float32 = 60
		_, _ = device.UpdateDemand(context.Background(), &electric.UpdateDemandRequest{
			Demand: &traits.ElectricDemand{
				Rating:  rating,
				Voltage: &voltage,
				Current: float32(math.Round(rand.Float64()*40*100) / 100),
			},
		})
		createElectricModes(device, rating)
		// set the active mode to the default one we just created (normal mode)
		_, _ = device.ClearActiveMode(context.Background(), &traits.ClearActiveModeRequest{})
		settings.Add(name, electric.WrapMemorySettings(device))
		return electric.Wrap(device), nil
	}
	return server.Collection(devices, settings)
}

func createElectricModes(device *electric.MemoryDevice, rating float32) {
	_, _ = device.CreateMode(context.Background(), &electric.CreateModeRequest{
		Mode: &traits.ElectricMode{
			Title:  "Normal Operation",
			Normal: true,
			Segments: []*traits.ElectricMode_Segment{
				{Magnitude: rating},
			},
		},
	})
	_, _ = device.CreateMode(context.Background(), &electric.CreateModeRequest{
		Mode: &traits.ElectricMode{
			Title: "Eco",
			Segments: []*traits.ElectricMode_Segment{
				{Magnitude: rating / 2},
			},
		},
	})
	_, _ = device.CreateMode(context.Background(), &electric.CreateModeRequest{
		Mode: &traits.ElectricMode{
			Title: "Quick Boot",
			Segments: []*traits.ElectricMode_Segment{
				{Magnitude: rating * 1.3},
			},
		},
	})
}
