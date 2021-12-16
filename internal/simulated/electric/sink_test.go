package electric

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"log"
	"time"
)

func ExampleSink() {
	clk := clock.Real()
	mem := electric.NewMemory(clk)
	sink := NewSink(mem,
		WithClock(clk),
		WithRampDuration(100*time.Millisecond),
	)

	// create a new normal mode
	mode, err := mem.CreateMode(&traits.ElectricMode{
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
	log.Println("Mode id:", mode.Id)

	// activate the normal mode
	_, err = mem.ChangeToNormalMode()
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
