package scrub

import "time"

// Scrubber allows jumping to a point in time.
// Types that implement Scrubber typically have some internal state that is a function of time.
// This interface allows for the adjustment of that state.
type Scrubber interface {
	// Scrub changes the state of this object to how it would be at the given time.
	Scrub(t time.Time) error
}
