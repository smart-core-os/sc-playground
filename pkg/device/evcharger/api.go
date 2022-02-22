package evcharger

import (
	"context"
	"log"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-playground/pkg/device"
	"github.com/smart-core-os/sc-playground/pkg/device/evcharger/event"
	"github.com/smart-core-os/sc-playground/pkg/sim/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApiServer struct {
	UnimplementedEVChargerApiServer

	clock      clock.Clock
	dispatcher input.Dispatcher
}

func (api *ApiServer) PlugIn(ctx context.Context, request *PlugInRequest) (*PlugInResponse, error) {
	if api.dispatcher == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "No external input supported")
	}
	log.Printf("Request %v PlugIn", request.Name)
	now := api.clock.Now()
	done, err := api.dispatcher.Dispatch(ctx, now, convertPlugInRequest(request, now))
	if err != nil {
		return nil, err
	}
	defer done()
	return &PlugInResponse{}, nil
}

func convertPlugInRequest(req *PlugInRequest, now time.Time) event.PlugIn {
	modes := make([]*event.ChargeMode, len(req.Event.Modes))
	for i, mode := range req.Event.Modes {
		modes[i] = &event.ChargeMode{
			Id:          mode.Id,
			Title:       mode.Title,
			Description: mode.Description,
			Segments:    mode.Segments,
		}
	}
	return event.PlugIn{
		Named: device.Named{Target: req.Name},
		Event: device.Event{Created: now},
		Modes: modes,
		Full:  req.Event.Full,
		Level: req.Event.Level,
	}
}

func (api *ApiServer) ChargeStart(ctx context.Context, request *ChargeStartRequest) (*ChargeStartResponse, error) {
	if api.dispatcher == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "No external input supported")
	}
	log.Printf("Request %v ChargeStart", request.Name)
	now := api.clock.Now()
	done, err := api.dispatcher.Dispatch(ctx, now, event.ChargeStart{
		Event: device.Event{Created: now},
		Named: device.Named{Target: request.Name},
	})
	if err != nil {
		return nil, err
	}
	defer done()
	return &ChargeStartResponse{}, nil
}

func (api *ApiServer) Unplug(ctx context.Context, request *UnplugRequest) (*UnplugResponse, error) {
	if api.dispatcher == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "No external input supported")
	}
	log.Printf("Request %v Unplug", request.Name)
	now := api.clock.Now()
	done, err := api.dispatcher.Dispatch(ctx, now, event.Unplug{
		Event: device.Event{Created: now},
		Named: device.Named{Target: request.Name},
	})
	if err != nil {
		return nil, err
	}
	defer done()
	return &UnplugResponse{}, nil
}
