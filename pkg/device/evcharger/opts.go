package evcharger

import (
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-playground/pkg/sim/input"
)

var (
	DefaultRating  float32 = 60  // Amps
	DefaultVoltage float32 = 230 // Volts
)

type Opt func(*opt)
type opt struct {
	clock       clock.Clock
	rating      float32
	voltage     float32
	dedicatedTL bool
	dispatcher  input.Dispatcher
}

var DefaultOpts = []Opt{
	WithClock(clock.Real()),
	WithVoltage(DefaultVoltage),
	WithRating(DefaultRating),
}

func WithClock(clk clock.Clock) Opt {
	return func(o *opt) {
		o.clock = clk
	}
}

func WithRating(r float32) Opt {
	return func(o *opt) {
		o.rating = r
	}
}

func WithVoltage(v float32) Opt {
	return func(o *opt) {
		o.voltage = v
	}
}

// WithDedicatedTL instructs the device that the TL instance is just for that device, no filtering will need to be done.
func WithDedicatedTL() Opt {
	return func(o *opt) {
		o.dedicatedTL = true
	}
}

func WithInputDispatcher(dispatcher input.Dispatcher) Opt {
	return func(o *opt) {
		o.dispatcher = dispatcher
	}
}

func calcSetup(opts ...Opt) *opt {
	setup := &opt{}
	for _, opt := range DefaultOpts {
		opt(setup)
	}
	for _, opt := range opts {
		opt(setup)
	}
	return setup
}
