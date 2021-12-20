package dynamic

import (
	"context"
	"sync"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/time/clock"

	"github.com/smart-core-os/sc-playground/internal/util/broadcast"
	"github.com/smart-core-os/sc-playground/pkg/profile"

	"go.uber.org/zap"
)

type Float32Option func(value *Float32)

var DefaultFloat32Options = []Float32Option{
	WithUpdateInterval(DefaultUpdateInterval),
	WithClock(clock.Real()),
	WithLogger(zap.NewNop()), // logging output from a Float32 is not useful in most cases so discard by default
}

// Float32 is a float32 that can change over time.
type Float32 struct {
	variable *broadcast.Variable

	// configuration options
	clock    clock.Clock
	interval time.Duration
	logger   *zap.Logger

	// manages running interpolations and profiles
	stopInterp  context.CancelFunc
	stopProfile context.CancelFunc
	stopM       sync.RWMutex
}

func NewFloat32(initial float32, options ...Float32Option) *Float32 {
	v := &Float32{}

	for _, opt := range DefaultFloat32Options {
		opt(v)
	}
	for _, opt := range options {
		opt(v)
	}

	v.variable = broadcast.NewVariable(v.clock, initial)

	return v
}

// Stop will stop any interpolations or profiles that are currently running on this Float32.
func (f *Float32) Stop() {
	f.stopM.RLock()
	defer f.stopM.RUnlock()

	if f.stopInterp != nil {
		f.stopInterp()
	}
	if f.stopProfile != nil {
		f.stopProfile()
	}
}

// Get retrieves the current value contained in this Float32.
func (f *Float32) Get() float32 {
	return f.variable.Value().(float32)
}

// Set will change this Float32 to a constant level. Any interpolations or profiles currently running will be stopped.
func (f *Float32) Set(value float32) {
	f.logger.Debug("set value of Float32 to constant",
		zap.Float32("value", value))

	f.Stop()
	f.variable.Set(value)
}

// Listen returns a Listener that sends updates whenever the Float32 changes for any reason.
// The values in events will always be float32.
func (f *Float32) Listen() *broadcast.Listener {
	return f.variable.Listen()
}

// StartInterpolation will begin a real-time linear interpolation of the Float32, from the current value
// to the target, taking place over the duration d. Only one interpolation may be in progress at a time.
// Any other interpolations or profiles that are in progress will be stopped, and the new interpolation run in
// their place.
// Cancel ctx to halt the interpolation - the Float32 will freeze as-is.
// The returned context (which inherits from the provided ctx) will be cancelled when the interpolation stops.
func (f *Float32) StartInterpolation(ctx context.Context, target float32, d time.Duration) Completion {
	return f.startInterpolation(ctx, target, d, true)
}

func (f *Float32) startInterpolation(ctx context.Context, target float32, d time.Duration, stopProfile bool) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	// cancel any profile currently executing (and optionally a profile as well)
	func() {
		f.stopM.Lock()
		defer f.stopM.Unlock()

		if stopProfile && f.stopProfile != nil {
			f.stopProfile()
		}
		if f.stopInterp != nil {
			f.stopInterp()
		}
		// record the new context's cancel function to allow stopping the new interpolation
		f.stopInterp = cancel
	}()

	startPoint := f.variable.Value().(float32)
	startTime := f.clock.Now()

	go func() {
		every := f.clock.Every(f.interval)
		defer every.Stop()
		defer cancel()

		for {
			select {
			case now := <-every.C():
				runDuration := now.Sub(startTime)
				proportion := float32(runDuration) / float32(d)
				if proportion > 1 {
					proportion = 1
				}
				newValue := startPoint + proportion*(target-startPoint)

				f.logger.Debug("interpolated Float32",
					zap.Float32("proportion", proportion),
					zap.Float32("value", newValue))

				f.variable.Set(newValue)

				if proportion >= 1 {
					// cancel once we reach the end
					f.logger.Debug("interpolation complete",
						zap.Duration("runDuration", runDuration))
					return
				}
			case <-ctx.Done():
				f.logger.Debug("context cancelled",
					zap.Error(ctx.Err()))
				return
			}
		}
	}()

	return ctx
}

// StartProfile will simulate a profile on the Float32, starting at the current time.
// The Float32 will be interpolated to each dynamic.Segment level in turn. The interpolation time is specified
// by the parameter change. After all segments are done, the value will be interpolated to profile.FinalLevel
// Cancel ctx to stop the profile from running. The value will be frozen as-is.
// The returned context.Context (which inherits from ctx) will be cancelled once the profile has finished simulating
// and the Float32 has reached profile.FinalLevel
func (f *Float32) StartProfile(ctx context.Context, profile profile.Profile, change time.Duration) Completion {
	ctx, cancel := context.WithCancel(ctx)
	func() {
		f.stopM.Lock()
		defer f.stopM.Unlock()

		if f.stopProfile != nil {
			f.stopProfile()
		}
		f.stopProfile = cancel
	}()

	startTime := f.clock.Now()

	go func() {
		defer cancel()
		segments := profile.Segments
		nextTime := startTime

		for {
			select {
			case <-f.clock.At(nextTime):
			case <-ctx.Done():
				f.logger.Debug("profile cancelled", zap.Error(ctx.Err()))
				return
			}

			if len(segments) == 0 {
				break
			}

			segment := segments[0]
			c := change
			if c > segment.Duration {
				c = segment.Duration
			}

			f.logger.Debug("starting new segment",
				zap.Float32("level", segment.Level),
				zap.Duration("duration", segment.Duration),
				zap.Duration("change", c),
			)

			_ = f.startInterpolation(ctx, segment.Level, c, false)

			nextTime = nextTime.Add(segment.Duration)
			segments = segments[1:]
		}

		// interpolate to final value, wait for it to finish
		f.logger.Debug("switching to FinalLevel",
			zap.Float32("level", profile.FinalLevel),
			zap.Duration("change", change),
		)
		<-f.startInterpolation(ctx, profile.FinalLevel, change, false).Done()
		f.logger.Debug("profile complete")
	}()

	return ctx

}

func WithClock(clk clock.Clock) Float32Option {
	return func(f *Float32) {
		f.clock = clk
	}
}

func WithUpdateInterval(interval time.Duration) Float32Option {
	return func(f *Float32) {
		f.interval = interval
	}
}

func WithLogger(logger *zap.Logger) Float32Option {
	return func(f *Float32) {
		f.logger = logger
	}
}
