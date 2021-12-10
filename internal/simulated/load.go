package simulated

import (
	"context"
	"github.com/smart-core-os/sc-playground/internal/simulated/demand"
	"github.com/smart-core-os/sc-playground/internal/util/broadcast"
	"github.com/smart-core-os/sc-playground/internal/util/clock"
	"sync"
	"time"
)

type Load struct {
	load  *broadcast.Variable
	clock clock.Clock

	stopInterp  context.CancelFunc
	stopProfile context.CancelFunc
	stopM       sync.RWMutex
}

func NewLoad(c clock.Clock) *Load {
	return &Load{
		load:  broadcast.NewVariable(c, float32(0)),
		clock: c,
	}
}

// Stop will stop any interpolations or profiles that are currently running on this Load.
func (l *Load) Stop() {
	l.stopM.RLock()
	defer l.stopM.RUnlock()

	if l.stopInterp != nil {
		l.stopInterp()
	}
	if l.stopProfile != nil {
		l.stopProfile()
	}
}

// Get retrieves the current value of this load.
func (l *Load) Get() float32 {
	return l.load.Value().(float32)
}

// Set will change this load to a constant level. Any interpolations or profiles currently running will be stopped.
func (l *Load) Set(value float32) {
	l.Stop()
	l.load.Set(value)
}

// Listen returns a Listener that sends updates whenever the load changes for any reason.
func (l *Load) Listen() *broadcast.Listener {
	return l.load.Listen()
}

// StartInterpolation will begin a real-time linear interpolation of the load value, from the current value
// to the target, taking place over the duration d. Only one interpolation
func (l *Load) StartInterpolation(ctx context.Context, target float32, d time.Duration) context.Context {
	return l.startInterpolation(ctx, target, d, true)
}

func (l *Load) startInterpolation(ctx context.Context, target float32, d time.Duration, stopProfile bool) context.Context {
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

	startPoint := l.load.Value().(float32)
	startTime := l.clock.Now()

	go func() {
		every := l.clock.Every(UpdateRate)
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

				l.load.Set(newLoad)

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

func (l *Load) StartProfile(ctx context.Context, profile demand.Profile, change time.Duration) context.Context {
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
