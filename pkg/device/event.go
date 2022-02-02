package device

import "time"

// Event is the basic type all device events should implement.
// Event is an timeline.Ater and can be placed onto a timeline.TL.
type Event struct {
	Created time.Time // Time the event was created
}

// At implements the timeline.Ater interface.
// It returns the Created time.
func (e Event) At() time.Time {
	return e.Created
}

// Named can be composed in an Event to allow it to be filtered by name.
// Named implements tlutil.Namer.
// See tlutil.FilterByName.
type Named struct {
	Target string
}

func (n Named) Name() string {
	return n.Target
}
