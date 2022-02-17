package event

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-playground/pkg/device"
	"github.com/smart-core-os/sc-playground/pkg/sim/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ElectricApiServer struct {
	traits.ElectricApiServer

	clock      clock.Clock
	dispatcher input.Dispatcher
}

func CaptureElectricInput(api traits.ElectricApiServer, clock clock.Clock, dispatcher input.Dispatcher) traits.ElectricApiServer {
	return &ElectricApiServer{
		ElectricApiServer: api,
		clock:             clock,
		dispatcher:        dispatcher,
	}
}

func (s *ElectricApiServer) Unwrap() interface{} {
	return s.ElectricApiServer
}

func (s *ElectricApiServer) ClearActiveMode(ctx context.Context, req *traits.ClearActiveModeRequest) (*traits.ElectricMode, error) {
	return s.setActiveMode(ctx, req.Name, "")
}

func (s *ElectricApiServer) UpdateActiveMode(ctx context.Context, req *traits.UpdateActiveModeRequest) (*traits.ElectricMode, error) {
	if req.ActiveMode.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ActiveMode.id should be set on update")
	}
	return s.setActiveMode(ctx, req.Name, req.ActiveMode.Id)
}

func (s *ElectricApiServer) setActiveMode(ctx context.Context, name, id string) (*traits.ElectricMode, error) {
	now := s.clock.Now()
	done, err := s.dispatcher.Dispatch(ctx, now, ModeChange{
		Named: device.Named{Target: name},
		Event: device.Event{Created: now},
		Id:    id, // blank id, means clear
	})
	if err != nil {
		return nil, err
	}
	defer done()
	return s.ElectricApiServer.GetActiveMode(ctx, &traits.GetActiveModeRequest{Name: name})
}
