package timeline

import "time"

// BoundTL is a TL with a Bound method.
type BoundTL interface {
	TL
	// Bound returns the first and last event time, if they exist.
	Bound() (first, last time.Time, exists bool)
}

// Bound returns the first and last event time of tl, if they exist.
//
// If tl implements BoundTl, Bound calls tl.Bound.
// Otherwise Bound attempts to find the first and last events by calling Previous and Next using MinTime and MaxTime.
func Bound(tl TL) (first, last time.Time, exists bool) {
	if b, ok := tl.(BoundTL); ok {
		return b.Bound()
	}

	first, exists = min(tl)
	if exists {
		last, _ = max(tl)
	}
	return first, last, exists
}

func min(tl TL) (min time.Time, exists bool) {
	if es := tl.At(MinTime); len(es) > 0 {
		return MinTime, true
	}
	return tl.Next(MinTime)
}

func max(tl TL) (max time.Time, exists bool) {
	if es := tl.At(MaxTime); len(es) > 0 {
		return MaxTime, true
	}
	return tl.Previous(MaxTime)
}
