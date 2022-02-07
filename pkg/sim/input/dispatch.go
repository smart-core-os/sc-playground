package input

import (
	"context"
	"time"
)

// Dispatcher allows callers to contribute input events to a timeline.
// Dispatch accepts the event to place on the timeline at the time returned by timeFunc.
// Dispatch blocks until the event has been applied to the simulation model, i.e. Scrub is now after the event time.
// The returned func must be called to notify that the caller is done using the model state in its current incarnation.
type Dispatcher interface {
	Dispatch(ctx context.Context, t time.Time, e ...interface{}) (done func(), err error)
}
