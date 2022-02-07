package timeline

import (
	"time"
)

func ExampleAddTL() {
	tl := writableTL()
	if wtl, ok := tl.(AddTL); ok {
		wtl.Add(time.Now(), entries...)
	}
}

func ExampleRemoveTL() {
	tl := writableTL()
	if wtl, ok := tl.(RemoveTL); ok {
		wtl.Remove(time.Unix(1000, 0), entries...)
	}
}

var (
	entries = []interface{}{
		time.Unix(100, 0),
		time.Unix(200, 0),
		time.Unix(300, 0),
	}
)

func writableTL() TL {
	return nil
}
