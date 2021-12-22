package dynamic

import (
	"context"
	"testing"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/time/clock"

	"github.com/smart-core-os/sc-playground/pkg/profile"

	"go.uber.org/zap"
)

func TestFloat32_StartInterpolation(t *testing.T) {
	clk := clock.Real()

	f := NewFloat32(0,
		WithClock(clk),
		WithLogger(zap.NewExample()),
	)

	// start interpolation
	var (
		target   float32 = 1
		duration         = 500 * time.Millisecond
	)
	complete := f.StartInterpolation(context.Background(), target, duration)
	// wait until interpolation complete
	select {
	case <-time.After(duration + 10*time.Millisecond):
		t.Fatal("test timed out")
	case <-complete.Done():
	}

	// verify that value is the target value
	var (
		expect = target
		actual = f.Get()
	)
	if actual != expect {
		t.Error(actual)
	}
}

func TestFloat32_StartProfile(t *testing.T) {
	clk := clock.Real()
	f := NewFloat32(0,
		WithClock(clk),
		WithLogger(zap.NewExample()),
	)

	prof := profile.Profile{
		Segments: []profile.Segment{
			{1 * time.Second, 1},
		},
		FinalLevel: 2,
	}

	complete := f.StartProfile(context.Background(), prof, 250*time.Millisecond)

	// before simulation run, value should be unchanged from initial value
	if value := f.Get(); value != 0 {
		t.Errorf("initial value not held: f.Get() = %f", value)
	}

	// run simulation for half a second puts us square in the middle of the segment
	time.Sleep(500 * time.Millisecond)

	// we should have the profile level by now
	if value := f.Get(); value != 1 {
		t.Errorf("not at expected profile level: f.Get() = %f", value)
	}

	// run simulation until the profile is done
	select {
	case <-time.After(1 * time.Second):
		t.Fatal("test timed out")
	case <-complete.Done():
	}

	// now we should be at FinalLevel
	if value := f.Get(); value != 2 {
		t.Errorf("did not reach FinalLevel: f.Get() = %f", value)
	}
}
