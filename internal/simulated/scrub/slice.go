package scrub

import (
	"time"
)

// Slice implements the Scrubber interface calling Scrub on each item in the slice.
type Slice []Scrubber

func (g Slice) Scrub(t time.Time) error {
	var rErr error
	for _, scrubber := range g {
		err := scrubber.Scrub(t)
		if err != nil && rErr == nil {
			rErr = err
		}
	}
	return rErr
}
