package input

import (
	"context"

	"github.com/smart-core-os/sc-playground/pkg/timeline"
)

// Capturer captures unprocessed input events onto a timeline.
// Capture returns a Session that should be either committed or rejected to allow event dispatchers to continue their
// processing.
type Capturer interface {
	Capture(ctx context.Context, tl timeline.AddTL) (Session, error)
}

// Session represents an active input processing session.
// Event dispatchers publish events but block until those events are captured and processed by the simulation loop.
// The Session provides the interface for the loop to notify the dispatcher whether processing succeeded or failed.
type Session interface {
	// Commit signals that event processing was successful.
	Commit()
	// Reject signals that event processing failed with error err.
	Reject(err error)
}
