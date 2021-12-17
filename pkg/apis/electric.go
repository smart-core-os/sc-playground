package apis

import (
	"context"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	simelectric "github.com/smart-core-os/sc-playground/internal/simulated/electric"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log"
	"time"

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
		mem := electric.NewMemory(clock.Real())
		device := electric.NewMemoryDevice(mem)

		// assign voltage and rating
		var voltage float32 = 240
		var rating float32 = 60
		_, err := mem.UpdateDemand(
			&traits.ElectricDemand{
				Voltage: &voltage,
				Rating:  rating,
			},
			&fieldmaskpb.FieldMask{Paths: []string{"voltage", "rating"}},
		)
		if err != nil {
			log.Printf("error assigning voltage & rating to new device %s: %v", name, err)
		}

		createElectricModes(mem, rating)

		// set the active mode to the default one we just created (normal mode)
		_, err = mem.ChangeToNormalMode()
		if err != nil {
			log.Printf("error changing to the normal mode on new device %s: %v", name, err)
		}

		// start the simulation
		go func() {
			sink := simelectric.NewSink(mem)
			err := sink.Simulate(context.Background())
			if err != nil {
				log.Printf("electric device %s stopped simulating with an error: %v", name, err)
			}
		}()

		settings.Add(name, electric.WrapMemorySettings(device))
		return electric.Wrap(device), nil
	}
	return server.Collection(devices, settings)
}

func createElectricModes(device *electric.Memory, rating float32) {
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
