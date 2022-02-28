package electric

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/smart-core-os/sc-playground/internal/simulated/dynamic"
	"github.com/smart-core-os/sc-playground/internal/util/broadcast"
)

type SourceOption func(source *Source)

func WithSourceLogger(logger *zap.Logger) SourceOption {
	return func(source *Source) {
		source.log = logger
	}
}

func WithSourceClock(clk clock.Clock) SourceOption {
	return func(source *Source) {
		source.clock = clk
	}
}

var DefaultSourceOptions = []SourceOption{
	WithSourceClock(clock.Real()),
	WithSourceLogger(zap.L()),
}

// Source is a simulation wrapper for an electric device that reports aggregate demand for multiple electric.Models.
// It simulated a device such as a current clamp sensor, that measures the current consumed by an entire circuit.
type Source struct {
	// Model is the datastore for this device's electric implementation. There is nothing to prevent this Model from
	// having electric modes, but the Source will not respond in any way to mode changes, so it's not useful.
	// Do not write to (change the pointer) this field; it is set once only in the constructor.
	Model *electric.Model
	// Extra represents load on this source that doesn't come from a downstream Model. This is useful to represent
	// non-Smart Core devices and their impact on the electric system.
	// Do not write to (change the pointer) this field; it is set once only in the constructor.
	Extra *dynamic.Float32
	// Downstream is a map of device names to electric.Model.
	Downstream *broadcast.Map

	clock clock.Clock
	log   *zap.Logger
}

func NewSource(model *electric.Model, options ...SourceOption) *Source {
	s := &Source{
		Model: model,
		Extra: dynamic.NewFloat32(0),
	}

	for _, opt := range DefaultSourceOptions {
		opt(s)
	}
	for _, opt := range options {
		opt(s)
	}

	s.Downstream = broadcast.NewMap(s.clock)

	return s
}

func (s *Source) Simulate(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	demandMask, err := fieldmaskpb.New(&traits.ElectricDemand{}, "current")
	if err != nil {
		panic(err) // not possible unless structure changes
	}

	demandChanges := make(chan demandChange)

	// goroutine handles updating the Model demand when demand changes
	group.Go(func() error {
		// holds the ElectricDemand.Current for each downstream Model
		demands := make(map[string]float32)
		// holds the value of s.Extra.Get()
		extraCurrent, extraListener := s.Extra.GetAndListen()
		defer extraListener.Close()

		for {
			select {
			case <-ctx.Done():
				// context done
				return ctx.Err()
			case change := <-demandChanges:
				// current on one of the downstream models changed
				demands[change.name] = change.current
			case event := <-extraListener.C:
				// s.Extra changed
				extraCurrent = event.NewValue.(float32)
			}

			// recalculate the total demand
			totalCurrent := extraCurrent
			for _, current := range demands {
				totalCurrent += current
			}

			// update the model with the new demand
			_, err := s.Model.UpdateDemand(
				&traits.ElectricDemand{Current: totalCurrent},
				resource.WithUpdateMask(demandMask),
			)
			if err != nil {
				return err
			}
			s.log.Debug("updated aggregated current", zap.Float32("current", totalCurrent))
		}
	})

	// goroutine creates workers for each downstream model that forward changes to demandChanges
	group.Go(func() error {
		initial, listener := s.Downstream.GetMembersAndListen()
		defer listener.Close()

		// map of device names to cleanup functions
		doneFuncs := make(map[string]func())

		// create workers for all the initial downstream models
		for name, model := range initial {
			model := model.(*electric.Model)
			doneFuncs[name] = s.spawnDownstreamWorker(ctx, group, demandChanges, name, model)
			s.log.Debug("registered new downstream model (initial)",
				zap.String("name", name))
		}

		// listen for changes to the Map of Models
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case event := <-listener.C:
				switch event.Type {
				case broadcast.MapAddEvent, broadcast.MapReplaceEvent:
					done, ok := doneFuncs[event.Key]
					if ok {
						done()
						s.log.Debug("removed downstream model (replacement)",
							zap.String("name", event.Key))
					}

					model := event.NewValue.(*electric.Model)
					doneFuncs[event.Key] = s.spawnDownstreamWorker(ctx, group, demandChanges, event.Key, model)
					s.log.Debug("registered new downstream model (update)",
						zap.String("name", event.Key))

				case broadcast.MapRemoveEvent:
					done, ok := doneFuncs[event.Key]
					if ok {
						done()
						delete(doneFuncs, event.Key)
						s.log.Debug("removed downstream model",
							zap.String("name", event.Key))
					} else {
						s.log.Warn("request to remove a model that is not registered",
							zap.String("name", event.Key))
					}

				default:
					s.log.Warn("unrecognised Map event",
						zap.Int("type", int(event.Type)))
				}
			}
		}
	})

	s.log.Debug("source simulation started")

	return group.Wait()
}

type demandChange struct {
	name    string
	current float32
}

func (s *Source) spawnDownstreamWorker(ctx context.Context, group *errgroup.Group, send chan<- demandChange,
	name string, model *electric.Model) (done func()) {

	demandMask, err := fieldmaskpb.New(&traits.ElectricDemand{}, "current")
	if err != nil {
		panic(err) // not possible unless structure changes
	}

	ctx, done = context.WithCancel(ctx)
	changes := model.PullDemand(ctx, resource.WithReadMask(demandMask))
	initial := model.Demand()

	s.log.Debug("sending initial current",
		zap.String("name", name),
		zap.Float32("current", initial.Current))
	send <- demandChange{
		name:    name,
		current: initial.Current,
	}

	group.Go(func() error {
		for change := range changes {
			s.log.Debug("sending current update",
				zap.String("name", name),
				zap.Float32("current", change.Value.Current))
			send <- demandChange{
				name:    name,
				current: change.Value.Current,
			}
		}
		return nil
	})

	return done
}
