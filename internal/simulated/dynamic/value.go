package dynamic

import (
	"context"
	"github.com/smart-core-os/sc-playground/internal/simulated"
	"github.com/smart-core-os/sc-playground/internal/util/broadcast"
	"github.com/smart-core-os/sc-playground/internal/util/clock"
	"sync"
	"time"
)

// Value is a float32 that can change over time.
type Value struct {
	value *broadcast.Variable
	clock clock.Clock

	stopInterp  context.CancelFunc
	stopProfile context.CancelFunc
	stopM       sync.RWMutex
}

func NewValue(initial float32, c clock.Clock) *Value {
	return &Value{
		value: broadcast.NewVariable(c, initial),
		clock: c,
	}
}

// Stop will stop any interpolations or profiles that are currently running on this Value.
func (l *Value) Stop() {
	l.stopM.RLock()
	defer l.stopM.RUnlock()

	if l.stopInterp != nil {
		l.stopInterp()
	}
	if l.stopProfile != nil {
		l.stopProfile()
	}
}

// Get retrieves the current value of this Value.
func (l *Value) Get() float32 {
	return l.value.Value().(float32)
}

// Set will change this Value's value to a constant level. Any interpolations or profiles currently running will be stopped.
func (l *Value) Set(value float32) {
	l.Stop()
	l.value.Set(value)
}

// Listen returns a Listener that sends updates whenever the value changes for any reason.
// The values in events will always be float32.
func (l *Value) Listen() *broadcast.Listener {
	return l.value.Listen()
}

// StartInterpolation will begin a real-time linear interpolation of the value, from the current value
// to the target, taking place over the duration d. Only one interpolation may be in progress at a time.
// Any other interpolations or profiles that are in progress will be stopped, and the new interpolation run in
// their place.
// Cancel ctx to halt the interpolation - the value will freeze as-is.
// The returned context (which inherits from the provided ctx) will be cancelled when the interpolation stops.
func (l *Value) StartInterpolation(ctx context.Context, target float32, d time.Duration) context.Context {
	return l.startInterpolation(ctx, target, d, true)
}

func (l *Value) startInterpolation(ctx context.Context, target float32, d time.Duration, stopProfile bool) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	// cancel any profile currently executing (and optionally a profile as well)
	func() {
		l.stopM.Lock()
		defer l.stopM.Unlock()

		if stopProfile && l.stopProfile != nil {
			l.stopProfile()
		}
		if l.stopInterp != nil {
			l.stopInterp()
		}
		// record the new context's cancel function to allow stopping the new interpolation
		l.stopInterp = cancel
	}()

	startPoint := l.value.Value().(float32)
	startTime := l.clock.Now()

	go func() {
		every := l.clock.Every(simulated.UpdateRate)
		defer every.Stop()
		defer cancel()

		for {
			select {
			case <-every.C():
				runDuration := l.clock.Now().Sub(startTime)
				proportion := float32(runDuration) / float32(d)
				if proportion > 1 {
					proportion = 1
				}
				newLoad := startPoint + proportion*(target-startPoint)

				l.value.Set(newLoad)

				if runDuration > d {
					// cancel once we reach the end
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return ctx
}

// StartProfile will simulate a profile on the Value, starting at the current time.
// The Value's value will be interpolated to each dynamic.Segment level in turn. The interpolation time is specified
// by the parameter change. After all segments are done, the value will be interpolated to profile.FinalLevel
// Cancel ctx to stop the profile from running. The value will be frozen as-is.
// The returned context.Context (which inherits from ctx) will be cancelled once the profile has finished simulating
// and the Value's value has reached profile.FinalLevel
func (l *Value) StartProfile(ctx context.Context, profile Profile, change time.Duration) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	func() {
		l.stopM.Lock()
		defer l.stopM.Unlock()

		if l.stopProfile != nil {
			l.stopProfile()
		}
		l.stopProfile = cancel
	}()

	startTime := l.clock.Now()

	go func() {
		defer cancel()
		segments := profile.Segments
		nextTime := startTime

		for {
			if delay := nextTime.Sub(l.clock.Now()); delay > 0 {
				select {
				case <-l.clock.After(delay):
				case <-ctx.Done():
					return
				}
			}

			if len(segments) == 0 {
				break
			}

			segment := segments[0]
			c := change
			if c > segment.Duration {
				c = segment.Duration
			}

			_ = l.startInterpolation(ctx, segment.Level, c, false)

			nextTime = nextTime.Add(segment.Duration)
			segments = segments[1:]
		}

		// interpolate to final value, wait for it to finish
		<-l.startInterpolation(ctx, profile.FinalLevel, change, false).Done()
	}()

	return ctx

}
