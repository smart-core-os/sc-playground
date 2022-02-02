package tlutil

import (
	"time"

	"github.com/smart-core-os/sc-playground/pkg/timeline"
)

// Namer describes types with a Name.
type Namer interface {
	Name() string
}

// FilterByName returns a TL that only includes events that are a Namer and whose Name() == name.
func FilterByName(tl timeline.TL, name string) timeline.TL {
	return timeline.Filter(tl, func(t time.Time, e interface{}) bool {
		if n, ok := e.(Namer); ok {
			return n.Name() == name
		}
		return false
	})
}
