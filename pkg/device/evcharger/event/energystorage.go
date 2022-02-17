package event

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-playground/pkg/device"
	"github.com/smart-core-os/sc-playground/pkg/sim/input"
)

type EnergyStorageApiServer struct {
	traits.EnergyStorageApiServer

	clock      clock.Clock
	dispatcher input.Dispatcher
}

func CaptureEnergyStorageInput(api traits.EnergyStorageApiServer, clock clock.Clock, dispatcher input.Dispatcher) traits.EnergyStorageApiServer {
	return &EnergyStorageApiServer{
		EnergyStorageApiServer: api,
		clock:                  clock,
		dispatcher:             dispatcher,
	}
}

func (s EnergyStorageApiServer) Unwrap() interface{} {
	return s.EnergyStorageApiServer
}

func (s EnergyStorageApiServer) Charge(ctx context.Context, request *traits.ChargeRequest) (*traits.ChargeResponse, error) {
	var done func()
	var err error
	now := s.clock.Now()
	if request.Charge {
		done, err = s.dispatcher.Dispatch(ctx, now, ChargeStart{
			Event: device.Event{Created: now},
			Named: device.Named{Target: request.Name},
		})
	} else {
		// todo: stop charging
		return &traits.ChargeResponse{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer done()
	return &traits.ChargeResponse{}, nil
}
