package electric

import (
	"context"
	"errors"
	"fmt"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/smart-core-os/sc-playground/internal/simulated/dynamic"
	"github.com/smart-core-os/sc-playground/internal/util/broadcast"
	"github.com/smart-core-os/sc-playground/internal/util/clock"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log"
	"sync"
	"time"
)

const DefaultRampDuration = 5 * time.Second

// Sink is a simulation wrapper for an electric device that consumes power and doesn't distribute it
// to any other Smart Core electric devices (it sinks power). In other words, it represents the leaf nodes of the power
// distribution tree.
// The Sink does not store its own state - it is designed to be backed by an electric.MemoryDevice from sc-golang,
// but any device that implements the electric API and the electric memory settings API correctly could be used.
type Sink struct {
	api    traits.ElectricApiClient
	memory electric.MemorySettingsApiClient
	name   string
	load   *dynamic.Value

	createNormal sync.Once // used to ensure the default mode is created only once
	normalId     string
}

// NewSink constructs a Sink.
// The gRPC clients api and memory must control the same underlying device when provided with the same device name.
// The Sink is designed for use with an electric.MemoryDevice, which does this.
// The clock c must provide time that is consistent with any timestamps etc. returned by calls to the gRPC APIs.
// When interacting with a remote device not under the control of local simulation, you likely need to use
// clock.Real.
func NewSink(c clock.Clock, api traits.ElectricApiClient, memory electric.MemorySettingsApiClient, name string) *Sink {
	s := &Sink{
		api:    api,
		memory: memory,
		name:   name,
		load:   dynamic.NewValue(0, c),
	}

	return s
}

// SetNormalMode updates (and creates if necessary) the normal mode for the device.
// Each electric device is permitted to have at most one mode annotated as normal - this API modifies that mode.
func (s *Sink) SetNormalMode(ctx context.Context, mode DeviceMode) (*traits.ElectricMode, error) {
	var err error
	// create the default mode if it has not already been created
	s.createNormal.Do(func() {
		var created *traits.ElectricMode
		created, err = s.memory.CreateMode(ctx, &electric.CreateModeRequest{
			Name: s.name,
			Mode: &traits.ElectricMode{
				Title:       mode.Title,
				Description: mode.Description,
				Normal:      true,
			},
		})

		s.normalId = created.Id
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create a normal mode: %w", err)
	} else if s.normalId == "" {
		return nil, errors.New("the first call to SetNormalMode failed to create the mode; create a new Sink")
	}

	protoMode := mode.ToProto()
	protoMode.Id = s.normalId
	protoMode.Title = mode.Title
	protoMode.Description = mode.Description

	mask, err := fieldmaskpb.New(protoMode, "title", "description", "segments")
	if err != nil {
		panic(err) // should never happen
	}

	protoMode, err = s.memory.UpdateMode(ctx, &electric.UpdateModeRequest{
		Name:       s.name,
		Mode:       protoMode,
		UpdateMask: mask,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update mode %s: %w", s.normalId, err)
	}
	return protoMode, nil
}

// CreateMode will add a new mode to the device. The mode is a non-normal mode.
// Returns the populated protobuf traits.ElectricMode, which contains the new mode's ID.
func (s *Sink) CreateMode(ctx context.Context, mode DeviceMode) (*traits.ElectricMode, error) {
	req := &electric.CreateModeRequest{
		Name: s.name,
		Mode: mode.ToProto(),
	}

	created, err := s.memory.CreateMode(ctx, req)
	return created, err
}

// ChangeMode will switch to the mode with ID modeId. If that mode does not exist on the device, an error
// will result. Mode IDs are not guaranteed to be consistent from run to run.
func (s *Sink) ChangeMode(ctx context.Context, modeId string) (*traits.ElectricMode, error) {
	mode := &traits.ElectricMode{Id: modeId}
	mask, err := fieldmaskpb.New(mode, "id")
	if err != nil {
		panic(err) // should be impossible
	}

	req := &traits.UpdateActiveModeRequest{
		Name:       s.name,
		ActiveMode: mode,
		UpdateMask: mask,
	}

	mode, err = s.api.UpdateActiveMode(ctx, req)
	if err != nil {
		return nil, err
	}
	return mode, nil
}

// ClearMode will switch the device to its normal mode. It is an error to clear a device
// if its normal mode has not been set.
func (s *Sink) ClearMode(ctx context.Context) (*traits.ElectricMode, error) {
	mode, err := s.api.ClearActiveMode(ctx, &traits.ClearActiveModeRequest{
		Name: s.name,
	})
	return mode, err
}

// GetDemand gets the demand, in Amps, that this device is currently placing on the electricity network.
func (s *Sink) GetDemand() float32 {
	return s.load.Get()
}

// ListenDemand creates a listener that issues updates every time the underlying load changes.
// The values are always float32.
func (s *Sink) ListenDemand() *broadcast.Listener {
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

// runUpdateDemand spawns a worker goroutine that updates the Memory Device with the new demand
// whenever the simulated load value changes.
func (s *Sink) runUpdateDemand(ctx context.Context) error {
	listener := s.load.Listen()

	defer listener.Close()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case event := <-listener.C:
			update := &traits.ElectricDemand{
				Current: event.NewValue.(float32),
			}

			mask, err := fieldmaskpb.New(update, "current")
			if err != nil {
				panic(err)
			}

			request := &electric.UpdateDemandRequest{
				Name:       s.name,
				Demand:     update,
				UpdateMask: mask,
			}

			_, err = s.memory.UpdateDemand(ctx, request)
			if err != nil {
				return fmt.Errorf("failed to update demand on the memory device: %w", err)
			}
		}
	}
}

// runUpdateMode spawns a worker goroutine that changes the profile being simulated by the load
// based on the mode selected by the memory device.
// Runs until an error occurs or the context is cancelled.
func (s *Sink) runUpdateMode(ctx context.Context) error {
	// open a stream to get mode changes
	stream, err := s.api.PullActiveMode(ctx, &traits.PullActiveModeRequest{
		Name: s.name,
	})
	if err != nil {
		return fmt.Errorf("can't pull active modes: %w", err)
	}

	// get the initial mode
	res, err := s.api.GetActiveMode(ctx, &traits.GetActiveModeRequest{
		Name: s.name,
	})
	if err != nil {
		return fmt.Errorf("can't get initial mode: %w", err)
	}

	// start the simulation of the initial mode
	currentMode := res.Id
	err = s.simulateMode(ctx, res)
	if err != nil {
		return fmt.Errorf("can't start simulation: %w", err)
	}

	defer func() {
		if err := stream.CloseSend(); err != nil {
			log.Printf("error closing PullActiveMode stream for %s: %v", s.name, err)
		}
	}()

	for {
		event, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("error receiving from PullActiveMode stream: %w", err)
		}

		for _, change := range event.Changes {
			// filter on device name
			if change.Name != s.name {
				continue
			}

			newMode := change.ActiveMode
			// skip if we are already running this profile
			if newMode.Id == currentMode {
				continue
			}

			// simulate the new mode
			currentMode = newMode.Id
			err = s.simulateMode(ctx, newMode)
			if err != nil {
				return fmt.Errorf("can't simulate the new mode: %w", err)
			}
		}
	}
}

func (s *Sink) simulateMode(ctx context.Context, mode *traits.ElectricMode) error {
	profile := DeviceModeFromProto(mode).Profile
	_ = s.load.StartProfile(ctx, profile, DefaultRampDuration)
	return nil
}

// DeviceMode is equivalent to a traits.ElectricMode, but with irrelevant fields removed and
// using a dynamic.Profile.
type DeviceMode struct {
	Title       string
	Description string
	Profile     dynamic.Profile
}

func DeviceModeFromProto(mode *traits.ElectricMode) DeviceMode {
	return DeviceMode{
		Title:       mode.Title,
		Description: mode.Description,
		Profile:     dynamic.ProfileFromProto(mode),
	}
}

func (m *DeviceMode) ToProto() *traits.ElectricMode {
	return &traits.ElectricMode{
		Title:       m.Title,
		Description: m.Description,
		Segments:    m.Profile.ToProto(),
	}
}
