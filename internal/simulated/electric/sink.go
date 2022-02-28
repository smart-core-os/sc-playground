package electric

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"

	"github.com/smart-core-os/sc-playground/internal/simulated/dynamic"
	"github.com/smart-core-os/sc-playground/internal/util/broadcast"
	"github.com/smart-core-os/sc-playground/pkg/profile"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

const (
	DefaultRampDuration   = 1 * time.Second
	DefaultUpdateInterval = dynamic.DefaultUpdateInterval
)

type SinkOption func(sink *Sink)

var DefaultSinkOptions = []SinkOption{
	WithRampDuration(DefaultRampDuration),
	WithInitialLoad(0),
	WithClock(clock.Real()),
	WithUpdateInterval(DefaultUpdateInterval),
	WithLogger(zap.L()),
}

// Sink is a simulation wrapper for an electric memory device that consumes power and doesn't distribute it
// to any other Smart Core electric devices (it sinks power). In other words, it represents the leaf nodes of the power
// distribution tree.
// The Sink does not store its own state - it is backed by an electric.Model from sc-golang.
type Sink struct {
	Model *electric.Model

	// fields configured with options
	rampDuration   time.Duration
	clock          clock.Clock
	initialLoad    float32
	updateInterval time.Duration
	logger         *zap.Logger

	// keeps track of the load value during operation
	load *dynamic.Float32

	createNormal sync.Once // used to ensure the default mode is created only once
	normalId     string
}

// NewSink constructs a Sink.
// The gRPC clients api and memory must control the same underlying device when provided with the same device name.
// The Sink is designed for use with an electric.MemoryDevice, which does this.
// If a clock is configured with WithClock, its timestamps must be consistent with those returned by the gRPC clients.
// When interacting with a remote device not under the control of local simulation, you likely need to use
// clock.Real, which is the default.
func NewSink(model *electric.Model, options ...SinkOption) *Sink {

	s := &Sink{
		Model: model,
	}

	for _, opt := range DefaultSinkOptions {
		opt(s)
	}
	for _, opt := range options {
		opt(s)
	}

	s.load = dynamic.NewFloat32(s.initialLoad,
		dynamic.WithUpdateInterval(s.updateInterval),
		dynamic.WithClock(s.clock),
		dynamic.WithLogger(s.logger.Named("load")),
	)

	return s
}

// GetDemand gets the demand, in Amps, that this device is currently placing on the electricity network.
func (s *Sink) GetDemand() float32 {
	return s.load.Get()
}

// ListenDemand creates a listener that issues updates every time the underlying load changes.
// The values are always float32.
func (s *Sink) ListenDemand() *broadcast.VariableListener {
	return s.load.Listen()
}

// Simulate will run profiles in real time when the underlying memory device changes modes, and will
// update the memory device's demand accordingly. Blocks until the context is cancelled or an error occurs.
func (s *Sink) Simulate(ctx context.Context) error {
	var group *errgroup.Group
	group, ctx = errgroup.WithContext(ctx)

	group.Go(func() error {
		return s.runUpdateDemand(ctx)
	})

	group.Go(func() error {
		return s.runUpdateMode(ctx)
	})

	return group.Wait()
}

// runUpdateDemand spawns a worker goroutine that updates Model with the new demand
// whenever the simulated load value changes.
func (s *Sink) runUpdateDemand(ctx context.Context) error {
	listener := s.load.Listen()

	defer listener.Close()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case event, ok := <-listener.C:
			if !ok {
				panic("logic error: channel should never close before loop exits")
			}

			current := event.NewValue.(float32)
			s.logger.Debug("sink device load changed",
				zap.Float32("load", current))

			update := &traits.ElectricDemand{
				Current: current,
			}

			mask, err := fieldmaskpb.New(update, "current")
			if err != nil {
				panic(err)
			}

			_, err = s.Model.UpdateDemand(update, resource.WithUpdateMask(mask))
			if err != nil {
				return fmt.Errorf("failed to update demand on the memory device: %w", err)
			}

			s.logger.Debug("updated electric memory device demand", zap.Float32("demand", current))
		}
	}
}

// runUpdateMode spawns a worker goroutine that changes the profile being simulated by the load
// based on the mode selected by the memory device.
// Runs until an error occurs or the context is cancelled.
func (s *Sink) runUpdateMode(ctx context.Context) error {
	// open a stream to get mode changes
	stream := s.Model.PullActiveMode(ctx)

	// get the initial mode
	initialMode := s.Model.ActiveMode()

	// start the simulation of the initial mode
	err := s.simulateMode(ctx, initialMode)
	if err != nil {
		return fmt.Errorf("can't start simulation: %w", err)
	}

	activeId := initialMode.Id
	for change := range stream {
		newMode := change.ActiveMode
		// skip if we are already running this profile
		if newMode.Id == activeId {
			continue
		}

		// simulate the new mode
		activeId = newMode.Id
		err = s.simulateMode(ctx, newMode)
		if err != nil {
			return fmt.Errorf("can't simulate the new mode: %w", err)
		}
	}

	return ctx.Err()
}

func (s *Sink) simulateMode(ctx context.Context, mode *traits.ElectricMode) error {
	s.logger.Info("sink switching mode",
		zap.Reflect("mode", mode))

	prof := DeviceModeFromProto(mode).Profile
	_ = s.load.StartProfile(ctx, prof, s.rampDuration)
	return nil
}

func WithRampDuration(duration time.Duration) SinkOption {
	if duration < 0 {
		panic("duration is negative")
	}
	return func(sink *Sink) {
		sink.rampDuration = duration
	}
}

func WithInitialLoad(load float32) SinkOption {
	return func(sink *Sink) {
		sink.initialLoad = load
	}
}

func WithClock(clk clock.Clock) SinkOption {
	return func(sink *Sink) {
		sink.clock = clk
	}
}

func WithUpdateInterval(interval time.Duration) SinkOption {
	return func(sink *Sink) {
		sink.updateInterval = interval
	}
}

func WithLogger(logger *zap.Logger) SinkOption {
	return func(sink *Sink) {
		sink.logger = logger
	}
}

// DeviceMode is equivalent to a traits.ElectricMode, but with irrelevant fields removed and
// using a dynamic.Profile.
type DeviceMode struct {
	Title       string
	Description string
	Profile     profile.Profile
}

func DeviceModeFromProto(mode *traits.ElectricMode) DeviceMode {
	return DeviceMode{
		Title:       mode.Title,
		Description: mode.Description,
		Profile:     profile.FromProto(mode.Segments),
	}
}

func (m *DeviceMode) ToProto() *traits.ElectricMode {
	return &traits.ElectricMode{
		Title:       m.Title,
		Description: m.Description,
		Segments:    m.Profile.ToProto(),
	}
}
