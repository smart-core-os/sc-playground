package electric

import (
	"context"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	simelectric "github.com/smart-core-os/sc-playground/internal/simulated/electric"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

func Activate(n *node.Node) {
	logger := zap.NewExample()

	var group *Group // instantiated later
	settings := electric.NewMemorySettingsApiRouter()
	devices := electric.NewApiRouter(
		electric.WithElectricApiClientFactory(electricClientFactory()),
		router.WithOnChange(func(change router.Change) {
			registerClientFactory(change, group, settings, n, logger)
		}),
	)

	// combines all electric devices demand into one (well multiple as it supports max, min, average, etc)
	group = NewGroup(electric.WrapApi(devices))
	group.Announce(parentElectricDeviceName(n), n, devices)

	n.AddRouter(devices, settings)
	n.AddTraitFactory(trait.Electric, func(name string, _ proto.Message) error {
		_, err := devices.Get(name)
		return err
	})
	n.AddClientFactory(trait.Electric, func(conn *grpc.ClientConn) interface{} {
		return traits.NewElectricApiClient(conn)
	})
}

func parentElectricDeviceName(n *node.Node) string {
	return n.Name() + "/g/electric"
}

func electricClientFactory() func(name string) (traits.ElectricApiClient, error) {
	return func(name string) (traits.ElectricApiClient, error) {
		model := electric.NewModel(clock.Real())

		// seed with a random load
		var voltage float32 = 240
		var rating float32 = 60
		_, err := model.UpdateDemand(&traits.ElectricDemand{
			Rating:  rating,
			Voltage: &voltage,
			Current: float32(math.Round(rand.Float64()*40*100) / 100),
		})
		if err != nil {
			log.Printf("error assigning voltage & rating to new device %s: %v", name, err)
		}
		createElectricModes(model, rating)
		// set the active mode to the default one we just created (normal mode)
		_, err = model.ChangeToNormalMode()
		if err != nil {
			log.Printf("error changing to the normal mode on new device %s: %v", name, err)
		}

		return electric.WrapApi(electric.NewModelServer(model)), nil
	}
}

func registerClientFactory(change router.Change, group *Group, settings *electric.MemorySettingsApiRouter, n *node.Node, logger *zap.Logger) {
	name := change.Name
	// add all devices to the group, avoiding cycles with the group being added itself
	if !strings.HasPrefix(name, parentElectricDeviceName(n)) {
		group.Add(name)
		log.Printf("%v now contributes towards %v demand", name, parentElectricDeviceName(n))
	}

	if !change.Auto {
		return
	}
	model, ok := wrap.UnwrapFully(change.New).(*electric.Model)
	if !ok {
		return
	}
	log.Printf("ElectricApiClient(%v) auto-created", name)

	// start the simulation
	go func() {
		logger := logger.Named("sink").With(zap.String("name", name))
		sink := simelectric.NewSink(model, simelectric.WithLogger(logger))
		err := sink.Simulate(context.Background())
		if err != nil {
			log.Printf("electric device %s stopped simulating with an error: %v", name, err)
		}
	}()

	settings.Add(name, electric.WrapMemorySettingsApi(electric.NewModelServer(model)))
	n.Announce(name, node.HasTrait(trait.Electric))
}

func createElectricModes(device *electric.Model, rating float32) {
	_, _ = device.CreateMode(&traits.ElectricMode{
		Title:  "Normal Operation",
		Normal: true,
		Segments: []*traits.ElectricMode_Segment{
			{Magnitude: rating * 0.8},
		},
	})
	_, _ = device.CreateMode(&traits.ElectricMode{
		Title: "Eco",
		Segments: []*traits.ElectricMode_Segment{
			{Magnitude: rating * 0.2, Length: durationpb.New(60 * time.Second)},
			{Magnitude: rating * 0.7},
		},
	})
	_, _ = device.CreateMode(&traits.ElectricMode{
		Title: "Quick Boot",
		Segments: []*traits.ElectricMode_Segment{
			{Magnitude: rating * 0.4, Length: durationpb.New(30 * time.Second)},
			{Magnitude: rating, Length: durationpb.New(10 * time.Second)},
			{Magnitude: rating * 0.8},
		},
	})
	_, _ = device.CreateMode(&traits.ElectricMode{
		Title: "Standby",
		Segments: []*traits.ElectricMode_Segment{
			{Magnitude: rating * 0.1},
		},
	})
}
