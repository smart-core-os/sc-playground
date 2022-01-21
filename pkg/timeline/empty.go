package timeline

// EmptyTL is a TL with an Empty method.
type EmptyTL interface {
	TL
	// Empty return true if there are no event in the timeline.
	Empty() bool
}

// Empty returns true if there are no events in the timeline.
//
// If tl implements EmptyTL, Empty calls tl.Empty.
// Otherwise Empty uses Bound to check for existence.
func Empty(tl TL) bool {
	if t, ok := tl.(EmptyTL); ok {
		return t.Empty()
	}

	_, _, exists := Bound(tl)
	return !exists
}
