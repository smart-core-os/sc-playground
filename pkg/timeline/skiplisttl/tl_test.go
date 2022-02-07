package skiplisttl

import (
	"fmt"
	"testing"
	"time"

	"github.com/smart-core-os/sc-playground/pkg/timeline"
)

func ExampleTL() {
	t1 := time.Date(2022, 02, 03, 15, 40, 10, 0, time.UTC)
	t2 := time.Date(2021, 12, 10, 8, 20, 0, 0, time.UTC)

	tl := New()
	tl.Add(t1, "A-1", "A-2")
	tl.Add(t2, "B")

	start, end, _ := timeline.Bound(tl)
	fmt.Printf("TL Bounds [%v, %v)", start, end)
	// Output: TL Bounds [2021-12-10 08:20:00 +0000 UTC, 2022-02-03 15:40:10 +0000 UTC)
}

func TestTL_Empty(t *testing.T) {
	tl := New()
	if !tl.Empty() {
		t.Fatalf("Empty: want true, got false")
	}

	tl.Add(time.Unix(1000, 0), "foo")
	if tl.Empty() {
		t.Fatalf("Empty: want false, got true")
	}
}

func TestTL_At(t *testing.T) {
	tl := New()
	got := tl.At(time.Unix(1000, 0))
	if len(got) > 0 {
		t.Fatalf("When empty, expect len=0. Got %v", got)
	}

	tl.Add(time.Unix(1000, 0), "foo")
	got = tl.At(time.Unix(1000-1, 0))
	if len(got) > 0 {
		t.Fatalf("When just before, expect len=0. Got %v", got)
	}
	got = tl.At(time.Unix(1000+1, 0))
	if len(got) > 0 {
		t.Fatalf("When just after, expect len=0. Got %v", got)
	}
	got = tl.At(time.Unix(1000, 0))
	if len(got) != 1 {
		t.Fatalf("When at time, expect len=1. Got %v", got)
	}
}

func TestTL_Next(t *testing.T) {
	tl := New()
	if next, ok := tl.Next(time.Unix(1000, 0)); ok {
		t.Fatalf("Empty want no next, got ( %v, true )", next)
	}

	tl.Add(time.Unix(1000, 0), "foo")
	if next, ok := tl.Next(time.Unix(1000-1, 0)); !ok || !next.Equal(time.Unix(1000, 0)) {
		t.Fatalf("Just before, want ( %v %v ), got ( %v %v )", time.Unix(1000, 0), true, next, ok)
	}
	if next, ok := tl.Next(time.Unix(1000, 0)); ok {
		t.Fatalf("At time, want ( _ false ), got ( %v %v )", next, ok)
	}
	if next, ok := tl.Next(time.Unix(1000+1, 0)); ok {
		t.Fatalf("Just after, want ( _ false ), got ( %v %v )", next, ok)
	}
}

func TestTL_Previous(t *testing.T) {
	tl := New()
	if next, ok := tl.Previous(time.Unix(1000, 0)); ok {
		t.Fatalf("Empty want no next, got ( %v, true )", next)
	}

	tl.Add(time.Unix(1000, 0), "foo")
	if prev, ok := tl.Previous(time.Unix(1000+1, 0)); !ok || !prev.Equal(time.Unix(1000, 0)) {
		t.Fatalf("Just after, want ( %v %v ), got ( %v %v )", time.Unix(1000, 0), true, prev, ok)
	}
	if prev, ok := tl.Previous(time.Unix(1000, 0)); ok {
		t.Fatalf("At time, want ( _ false ), got ( %v %v )", prev, ok)
	}
	if prev, ok := tl.Previous(time.Unix(1000-1, 0)); ok {
		t.Fatalf("Just before, want ( _ false ), got ( %v %v )", prev, ok)
	}
}
