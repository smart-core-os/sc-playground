package sim

import (
	"context"
	"fmt"
	"time"

	"github.com/smart-core-os/sc-playground/pkg/sim/input"
	"github.com/smart-core-os/sc-playground/pkg/sim/scrub"
	"github.com/smart-core-os/sc-playground/pkg/sim/stats"
	"github.com/smart-core-os/sc-playground/pkg/timeline"
)

// Loop repeatedly advances state by recording events on a timeline and scrubbing ahead of those events.
type Loop struct {
	MinTimePerFrame time.Duration // how long a frame is, frames might take longer than this

	Input input.Capturer // captures user inputs
	TL    timeline.TL    // records user inputs and other info on a timeline
	Model scrub.Scrubber // the simulation model, the things we are controlling/exposing

	framer stats.Framer
}

func NewLoop(tl timeline.TL, model scrub.Scrubber, input input.Capturer, opts ...Option) *Loop {
	l := &Loop{
		MinTimePerFrame: time.Second / 60,
		Input:           input,
		TL:              tl,
		Model:           model,
	}
	applyOpts(l, opts...)
	return l
}

// Run executes the loop for as long as ctx is not cancelled or a critical error has not occurred.
func (l *Loop) Run(ctx context.Context) error {
	// Things we could do in this loop:
	// 1. Have a more dynamic update rate, if we're constantly dropping updates, switch to something less frequent.
	//    This is based on the assumption that a steady update rate is a desirable thing.
	// 2. Record statistics for how we're doing - actual ups, number dropped frames, idle time, etc
	// 3. Report stats via an API we can consume on the UI. I think this would be useful to show

	ticker := time.NewTicker(l.MinTimePerFrame)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-ticker.C:
			if err := l.Advance(ctx, t); err != nil {
				return err
			}
		}
	}
}

// Advance captures outstanding events and scrubs the model to the given time.
// There is no need to call Advance if you have called Run, it does it for you.
func (l *Loop) Advance(ctx context.Context, t time.Time) (err error) {
	var session input.Session
	t0 := time.Now()
	frame := stats.FrameTiming{}
	defer func() {
		frame.Total = time.Since(t0)
		if l.framer != nil {
			l.framer.Frame(frame)
		}
	}()
	defer func() {
		if err != nil && session != nil {
			t0 := time.Now()
			session.Reject(err)
			frame.Respond = time.Since(t0)
		}
	}()

	// 1. Capture "user" inputs, and update the timeline
	if atl, ok := l.TL.(timeline.AddTL); ok {
		t0 := time.Now()
		session, err = l.Input.Capture(ctx, atl)
		frame.Capture = time.Since(t0)
		if err != nil {
			err = fmt.Errorf("during input capture: %w", err)
			return
		}
	}

	// 2. Scrub - updates sim model to match the new state of the timeline at t
	s0 := time.Now()
	err = l.Model.Scrub(t)
	frame.Scrub = time.Since(s0)
	if err != nil {
		err = fmt.Errorf("during model scrub: %w", err)
		return
	}

	// 3. Let the input dispatchers know the state is correct, and wait for them to be done
	if session != nil {
		t0 := time.Now()
		session.Commit()
		frame.Respond = time.Since(t0)
		session = nil
	}

	// todo:
	//   4. Make Keyframes

	// todo:
	//   5. Clean up TL: trim, GC, etc. Make sure we aren't growing exponentially

	return
}
