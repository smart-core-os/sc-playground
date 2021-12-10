package simulated

import (
	"context"
	"errors"
	"fmt"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/smart-core-os/sc-playground/internal/simulated/demand"
	"github.com/smart-core-os/sc-playground/internal/util/broadcast"
	"github.com/smart-core-os/sc-playground/internal/util/clock"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log"
	"sync"
	"time"
)

const rampDuration = 5 * time.Second

type SinkDevice struct {
	api    traits.ElectricApiClient
	memory electric.MemorySettingsApiClient
	name   string
	load   *Load
	idle   *broadcast.Variable

	createNormal sync.Once // used to ensure the default mode is created only once
	normalId     string
}

func NewSinkDevice(c clock.Clock, api traits.ElectricApiClient, memory electric.MemorySettingsApiClient, name string) *SinkDevice {
	s := &SinkDevice{
		api:    api,
		memory: memory,
		name:   name,
		load:   NewLoad(c),
		idle:   broadcast.NewVariable(c, true),
	}

	return s
}

func (s *SinkDevice) SetNormalMode(ctx context.Context, mode ElectricDeviceMode) (*traits.ElectricMode, error) {
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
		return nil, errors.New("the first call to SetNormalMode failed to create the mode; create a new SinkDevice")
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

func (s *SinkDevice) ChangeMode(ctx context.Context, modeId string) (*traits.ElectricMode, error) {
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
func (s *SinkDevice) ClearMode(ctx context.Context) (*traits.ElectricMode, error) {
	mode, err := s.api.ClearActiveMode(ctx, &traits.ClearActiveModeRequest{
		Name: s.name,
	})
	return mode, err
}

// ListenDemand creates a listener that issues updates every time the underlying load changes.
// The values are always float32.
func (s *SinkDevice) ListenDemand() *broadcast.Listener {
	return s.load.Listen()
}

// ListenIdle creates a listener that issues updates every time the device goes into or out of idle mode.
// The device is considered idle if it has reached the final segment of its currently simulated electric mode,
// after which there will be no more changes to demand until a new mode is selected.
// The listener will send true when the device becomes idle, and false when it is no longer idle.
func (s *SinkDevice) ListenIdle() *broadcast.Listener {
	panic("not implemented")
}

// Simulate will run profiles in real time when the underlying memory device changes modes, and will
// update the memory device's demand accordingly. Blocks until the context is cancelled or an error occurs.
func (s *SinkDevice) Simulate(ctx context.Context) error {
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
func (s *SinkDevice) runUpdateDemand(ctx context.Context) error {
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
func (s *SinkDevice) runUpdateMode(ctx context.Context) error {
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

func (s *SinkDevice) simulateMode(ctx context.Context, mode *traits.ElectricMode) error {
	profile := ElectricDeviceModeFromProto(mode).Profile
	_ = s.load.StartProfile(ctx, profile, rampDuration)
	return nil
}

type ElectricDeviceMode struct {
	Title       string
	Description string
	Profile     demand.Profile
}

func ElectricDeviceModeFromProto(mode *traits.ElectricMode) ElectricDeviceMode {
	return ElectricDeviceMode{
		Title:       mode.Title,
		Description: mode.Description,
		Profile:     demand.ProfileFromProto(mode),
	}
}

func (m *ElectricDeviceMode) ToProto() *traits.ElectricMode {
	return &traits.ElectricMode{
		Title:       m.Title,
		Description: m.Description,
		Segments:    m.Profile.ToProto(),
	}
}
