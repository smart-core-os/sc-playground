// Package timeline defines the basic types of a timeline.
package timeline

import (
	"time"
)

// TL implements an ordered set of events, a timeline.
type TL interface {
	// At returns all the events at the specific time.
	At(t time.Time) []interface{}

	// Previous returns the newest time with events < t on this TL.
	Previous(t time.Time) (previous time.Time, exists bool)
	// Next returns the oldest time with events > t on this TL.
	Next(t time.Time) (next time.Time, exists bool)
}

// MinTime is a sensible time for which timeline operations make sense.
// Timelines compare times to order them, these values prevent overflows during these operations.
var MinTime = time.Time{}

// MaxTime is the largest time after MinTime where durations don't overflow.
var MaxTime = time.Unix(1<<63-1-unixToInternal, 999999999)

// number of seconds between Year 1 and 1970 (62135596800 seconds)
// see https://stackoverflow.com/a/32620397/317404
var unixToInternal = int64((1969*365 + 1969/4 - 1969/100 + 1969/400) * 24 * 60 * 60)
