package timeline

import (
	"time"
)

// AddTL is a TL with an Add method.
//
// Add does not have a fallback mechanism, so no Add utility function exists for this feature.
// To add to a TL, explicitly cast the TL to an AddTL and call Add directly.
type AddTL interface {
	TL
	// Add adds the given entries to this TL at the time specified by t.
	Add(t time.Time, entries ...interface{})
}

// RemoveTL is a TL with a Remove method.
//
// Remove does not have a fallback mechanism, so no Remove utility function exists for this feature.
// To remove from a TL, explicitly cast the TL to a RemoveTL and call Remove directly.
type RemoveTL interface {
	TL
	// Remove removes the entries from the TL at the time specified by t.
	Remove(t time.Time, entries ...interface{})
}

// WriteTL is the interface that groups AddTL and RemoveTL
type WriteTL interface {
	AddTL
	RemoveTL
}
