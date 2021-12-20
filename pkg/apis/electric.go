package apis

import (
	"context"
	"log"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/memory"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	simelectric "github.com/smart-core-os/sc-playground/internal/simulated/electric"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"google.golang.org/protobuf/types/known/durationpb"
)

func ElectricApi() server.GrpcApi {
	devices := electric.NewRouter()
	settings := electric.NewMemorySettingsRouter()

	mask, err := fieldmaskpb.New(&traits.ElectricDemand{}, "voltage", "rating")
	if err != nil {
		panic(err) // should be impossible unless field names changed
	}

	devices.Factory = func(name string) (traits.ElectricApiClient, error) {
		log.Printf("Creating ElectricClient(%v)", name)
		model := electric.NewModel(clock.Real())

		// assign voltage and rating
		var voltage float32 = 240
		var rating float32 = 60
		_, err := model.UpdateDemand(
			&traits.ElectricDemand{
				Voltage: &voltage,
				Rating:  rating,
			},
			memory.WithUpdateMask(mask),
		)
		if err != nil {
			log.Printf("error assigning voltage & rating to new device %s: %v", name, err)
		}

		createElectricModes(model, rating)

		// set the active mode to the default one we just created (normal mode)
		_, err = model.ChangeToNormalMode()
		if err != nil {
			log.Printf("error changing to the normal mode on new device %s: %v", name, err)
		}

		// start the simulation
		go func() {
			sink := simelectric.NewSink(model)
			err := sink.Simulate(context.Background())
			if err != nil {
				log.Printf("electric device %s stopped simulating with an error: %v", name, err)
			}
		}()

		electricServer := electric.NewModelServer(model)
		settings.Add(name, electric.WrapMemorySettings(electricServer))
		return electric.Wrap(electricServer), nil
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
	},
	)
	_, _ = device.CreateMode(&traits.ElectricMode{
		Title: "Eco",
		Segments: []*traits.ElectricMode_Segment{
			{Magnitude: rating * 0.3, Length: durationpb.New(60 * time.Second)},
			{Magnitude: rating * 0.9, Length: durationpb.New(10 * time.Second)},
			{Magnitude: rating * 0.8},
		},
	},
	)
	_, _ = device.CreateMode(&traits.ElectricMode{
		Title: "Quick Boot",
		Segments: []*traits.ElectricMode_Segment{
			{Magnitude: rating * 1.3, Length: durationpb.New(30 * time.Second)},
			{Magnitude: rating},
		},
	},
	)
}
