package electric

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/smart-core-os/sc-playground/internal/simulated"
	"github.com/smart-core-os/sc-playground/internal/simulated/dynamic"
	"time"
)

func ExampleSink() {
	// create device
	dev := electric.NewMemoryDevice()
	api := electric.Wrap(dev)
	mem := electric.WrapMemorySettings(dev)
	clk := simulated.NewClock(time.Now())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect sink to control the device and create a mode
	sink := NewSink(api, mem, "ELEC-001",
		WithClock(clk),
		WithRampDuration(100*time.Millisecond))

	mode, err := sink.CreateMode(ctx, DeviceMode{
		Title:       "On",
		Description: "Device is powered on",
		Profile: dynamic.Profile{
			FinalLevel: 10,
		},
	})
	if err != nil {
		panic(err)
	}

	// select the new mode
	_, err = sink.ChangeMode(ctx, mode.Id)
	if err != nil {
		panic(err)
	}

	go func() {
		err := sink.Simulate(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	// wait for the mode to take effect
	simulated.SimulateFor(clk, 1*time.Second, 100*time.Millisecond)
	time.Sleep(1 * time.Second)

	fmt.Println(sink.GetDemand())

	// Output: 10
}
