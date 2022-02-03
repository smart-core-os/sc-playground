package skiplisttl

import (
	"fmt"
	"testing"

	"github.com/MauriceGit/skiplist"
	"github.com/google/go-cmp/cmp"
)

func TestFixedSkipList_Delete(t *testing.T) {
	sl := makeSkipList(
		element{1, "A"},
		element{1, "B"},
		element{1, "C"},
	)

	sl.Delete(element{1, "B"})

	want := makeSkipList(
		element{1, "A"},
		element{1, "C"},
	)

	assertSkipListEqual(t, want, sl)
}

func TestFixedSkipList_DeleteAll(t *testing.T) {
	sl := makeSkipList(
		element{1, "A"},
		element{1, "B"},
		element{1, "C"},
		element{2, "D"},
	)

	sl.DeleteAll(keyEntry(1))

	want := makeSkipList(element{2, "D"})

	assertSkipListEqual(t, want, sl)
}

func TestFixedSkipList_FindAll(t *testing.T) {
	sl := makeSkipList(
		element{1, "A"},
		element{1, "B"},
		element{1, "C"},
		element{2, "D"},
	)

	found := sl.FindAll(keyEntry(1))
	want := []skiplist.ListElement{
		element{1, "A"},
		element{1, "B"},
		element{1, "C"},
	}
	got := make([]skiplist.ListElement, len(found))
	for i, e := range found {
		got[i] = e.GetValue()
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("FindAll (-want, +got)\n%v", diff)
	}
}

func TestFixedSkipList_Bounds(t *testing.T) {
	sl := makeSkipList(
		element{1, "A"},
		element{1, "B"},
		element{1, "C"},
	)

	smallest, largest := sl.GetSmallestNode(), sl.GetLargestNode()
	if smallest.GetValue() != (element{1, "A"}) {
		t.Fatalf("Smallest want %v, got %v", element{1, "A"}, smallest)
	}
	if largest.GetValue() != (element{1, "C"}) {
		t.Fatalf("Largest want %v, got %v", element{1, "C"}, smallest)
	}
}

func makeSkipList(elements ...element) *fixedSkipList {
	underlying := skiplist.NewSeed(0)
	sl := &fixedSkipList{&underlying}
	for _, e := range elements {
		sl.Insert(e)
	}
	return sl
}

func assertSkipListEqual(t *testing.T, a, b *fixedSkipList) {
	if a == b {
		return
	}

	// why are they different?
	aFirst, bFirst := a.GetSmallestNode(), b.GetSmallestNode()
	if aFirst.GetValue() != bFirst.GetValue() {
		t.Fatalf("Element[0]: %v != %v", aFirst.GetValue(), bFirst.GetValue())
	}

	aNext, bNext := aFirst, bFirst
	for i := 1; ; i++ {
		aNext, bNext = a.Next(aNext), b.Next(bNext)
		if aNext == aFirst || bNext == bFirst {
			if aNext.GetValue() != bNext.GetValue() {
				if aNext == aFirst {
					t.Fatalf("Loop: [0]=%v != [%v]=%v", aNext.GetValue(), i, bNext.GetValue())
				} else {
					t.Fatalf("Loop: [%v]=%v != [0]=%v\n%v\n%v", i, aNext.GetValue(), bNext.GetValue(), a.String(), b.String())
				}
			}
			return
		}

		if aNext.GetValue() != bNext.GetValue() {
			t.Fatalf("Element[%v]: %v != %v", i, aFirst.GetValue(), bFirst.GetValue())
		}
	}
}

type element struct {
	Key   float64
	Value string
}

func (e element) ExtractKey() float64 {
	return e.Key
}

func (e element) String() string {
	return fmt.Sprintf("[%v]%v", e.Key, e.Value)
}
