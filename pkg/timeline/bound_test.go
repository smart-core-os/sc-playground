package timeline

import (
	"testing"
	"time"
)

type boundOnly struct{ BoundTL }

func (o boundOnly) At(t time.Time) []interface{} {
	return nil
}
func (o boundOnly) Previous(t time.Time) (previous time.Time, exists bool) {
	return
}
func (o boundOnly) Next(t time.Time) (next time.Time, exists bool) {
	return
}

func TestBound(t *testing.T) {
	check := func(name string, tl TL) {
		t.Helper()
		first, last, exists := Bound(tl)
		if !exists {
			t.Fatalf("%v bounds should exist", name)
		}
		if !first.Equal(testTL[0].At()) {
			t.Fatalf("%v first want %v, got %v", name, testTL[0].At(), first)
		}
		if !last.Equal(testTL[len(testTL)-1].At()) {
			t.Fatalf("%v last want %v, got %v", name, testTL[len(testTL)-1].At(), first)
		}
	}

	check("boundOnly", boundOnly{testTL})
	check("tlOnly", tlOnly{testTL})
}
