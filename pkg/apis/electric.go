package apis

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/time/clock"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"google.golang.org/protobuf/types/known/durationpb"
)

func ElectricApi() server.GrpcApi {
	devices := electric.NewRouter()
	settings := electric.NewMemorySettingsRouter()
	devices.Factory = func(name string) (traits.ElectricApiClient, error) {
		log.Printf("Creating ElectricClient(%v)", name)
		model := electric.NewModel(clock.Real())
		// seed with a random load
		var voltage float32 = 240
		var rating float32 = 60
		_, _ = model.UpdateDemand(&traits.ElectricDemand{
			Rating:  rating,
			Voltage: &voltage,
			Current: float32(math.Round(rand.Float64()*40*100) / 100),
		})
		createElectricModes(model, rating)
		// set the active mode to the default one we just created (normal mode)
		_, _ = model.ChangeToNormalMode()

		device := electric.NewModelServer(model)
		settings.Add(name, electric.WrapMemorySettings(device))
		return electric.Wrap(device), nil
	}
	return server.Collection(devices, settings)
}

func createElectricModes(device *electric.Model, rating float32) {
	_, _ = device.CreateMode(&traits.ElectricMode{
		Title:  "Normal Operation",
		Normal: true,
		Segments: []*traits.ElectricMode_Segment{
			{Magnitude: rating},
		},
	})
	_, _ = device.CreateMode(&traits.ElectricMode{
		Title: "Eco",
		Segments: []*traits.ElectricMode_Segment{
			{Magnitude: rating * 0.3, Length: durationpb.New(60 * time.Second)},
			{Magnitude: rating * 0.9, Length: durationpb.New(10 * time.Second)},
			{Magnitude: rating * 0.8},
		},
	})
	_, _ = device.CreateMode(&traits.ElectricMode{
		Title: "Quick Boot",
		Segments: []*traits.ElectricMode_Segment{
			{Magnitude: rating * 1.3, Length: durationpb.New(30 * time.Second)},
			{Magnitude: rating},
		},
	})
}
