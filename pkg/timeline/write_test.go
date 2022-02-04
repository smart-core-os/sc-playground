package timeline

import (
	"time"
)

func ExampleAddTL() {
	tl := writableTL()
	if wtl, ok := tl.(AddTL); ok {
		wtl.Add(timeFunc, entries...)
	}
}

func ExampleRemoveTL() {
	tl := writableTL()
	if wtl, ok := tl.(RemoveTL); ok {
		wtl.Remove(timeFunc, entries...)
	}
}

var (
	timeFunc = func(e interface{}) time.Time {
		return e.(time.Time)
	}
	entries = []interface{}{
		time.Unix(100, 0),
		time.Unix(200, 0),
		time.Unix(300, 0),
	}
)

func writableTL() TL {
	return nil
}
