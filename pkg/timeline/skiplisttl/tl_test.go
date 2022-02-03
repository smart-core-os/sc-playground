package skiplisttl

import (
	"fmt"
	"time"

	"github.com/smart-core-os/sc-playground/pkg/timeline"
)

func ExampleTL() {
	// Entry is the type we are storing in our TL.
	type Entry struct {
		at   time.Time
		data string
	}
	// timeFunc is how we convert an Entry to a time.Time.
	timeFunc := func(e interface{}) time.Time {
		return e.(Entry).at
	}

	t1 := time.Date(2022, 02, 03, 15, 40, 10, 0, time.UTC)
	t2 := time.Date(2021, 12, 10, 8, 20, 0, 0, time.UTC)

	tl := New()
	tl.Add(timeFunc,
		Entry{t1, "A-1"},
		Entry{t1, "A-2"},
		Entry{t2, "B"},
	)

	start, end, _ := timeline.Bound(tl)
	fmt.Printf("TL Bounds [%v, %v)", start, end)
	// Output: TL Bounds [2021-12-10 08:20:00 +0000 UTC, 2022-02-03 15:40:10 +0000 UTC)
}
