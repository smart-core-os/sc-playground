package dynamic

import (
	"context"
	"testing"
	"time"

	"github.com/smart-core-os/sc-playground/internal/simulated"
	profile2 "github.com/smart-core-os/sc-playground/pkg/profile"

	"go.uber.org/zap"
)

func TestFloat32_StartInterpolation(t *testing.T) {
	clk := simulated.NewClock(time.Now())

	f := NewFloat32(0,
		WithClock(clk),
		WithLogger(zap.NewExample()),
	)

	ctx := f.StartInterpolation(context.Background(), 1, 1*time.Second)

	simulated.SimulateFor(clk, 1100*time.Millisecond, 100*time.Millisecond)
	<-ctx.Done()

	var (
		expect float32 = 1.0
		actual         = f.Get()
	)
	if actual != expect {
		t.Error(actual)
	}
}

func TestFloat32_StartProfile(t *testing.T) {
	clk := simulated.NewClock(time.Now())
	f := NewFloat32(0,
		WithClock(clk),
		WithLogger(zap.NewExample()),
	)

	profile := profile2.Profile{
		Segments: []profile2.Segment{
			{time.Minute, 1},
		},
		FinalLevel: 2,
	}

	ctx := f.StartProfile(context.Background(), profile, time.Second)

	// before simulation run, value should be unchanged from initial value
	if value := f.Get(); value != 0 {
		t.Errorf("initial value not held: f.Get() = %f", value)
	}

	// run simulation for 30 seconds puts us square in the middle of the segment
	simulated.SimulateFor(clk, 30*time.Second, 100*time.Millisecond)

	// we should have the profile level by now
	if value := f.Get(); value != 1 {
		t.Errorf("not at expected profile level: f.Get() = %f", value)
	}

	// run simulation to the end
	simulated.SimulateFor(clk, 1*time.Minute, 100*time.Millisecond)

	<-ctx.Done()

	// now we should be at FinalLevel
	if value := f.Get(); value != 2 {
		t.Errorf("did not reach FinalLevel: f.Get() = %f", value)
	}
}
