package electric

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
)

func ExampleSink() {
	clk := clock.Real()
	model := electric.NewModel(clk)
	sink := NewSink(model,
		WithClock(clk),
		WithRampDuration(100*time.Millisecond),
	)

	// create a new mode
	mode, err := model.CreateMode(&traits.ElectricMode{
		Title:       "On",
		Description: "Device is powered on",
		Segments: []*traits.ElectricMode_Segment{
			{Magnitude: 10}, // maintain 10 amps indefinitely
		},
		Normal: true,
	})
	if err != nil {
		panic(err)
	}

	// activate the mode
	_, err = model.ChangeActiveMode(mode.Id)
	if err != nil {
		panic(err)
	}

	// run the simulation for 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = sink.Simulate(ctx)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		panic(err) // unexpected error in simulation
	}

	// log final current
	fmt.Println(sink.GetDemand())

	// Output: 10
}
