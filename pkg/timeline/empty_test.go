package timeline

import (
	"testing"
	"time"
)

type emptyOnly struct{ EmptyTL }

func (o emptyOnly) At(t time.Time) []interface{} {
	return nil
}
func (o emptyOnly) Previous(t time.Time) (previous time.Time, exists bool) {
	return
}
func (o emptyOnly) Next(t time.Time) (next time.Time, exists bool) {
	return
}

func TestEmpty(t *testing.T) {
	if Empty(emptyOnly{testTL}) {
		t.Fatalf("EmptyTL.Empty() should be false")
	}
	if Empty(tlOnly{testTL}) {
		t.Fatalf("Empty(TL) should be false")
	}

	empty := aterTL{}
	if !Empty(emptyOnly{empty}) {
		t.Fatalf("EmptyTL.Empty() should be true")
	}
	if !Empty(tlOnly{empty}) {
		t.Fatalf("Empty(TL) should be true")
	}
}
