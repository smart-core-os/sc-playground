// Package skiplisttl provides a skip list based implementation of timeline.TL.
package skiplisttl

import (
	"fmt"
	"time"

	"github.com/MauriceGit/skiplist"
)

// TL is an implementation of timeline.TL backed by a skip list.
// The TL is writable and any timeline.Slice or timeline.Filter applied to this TL will reflect changes, but will not
// be writable themselves.
type TL struct {
	items *fixedSkipList // of valueEntry
	key   KeyFunc
}

// New creates a new TL.
// See DefaultOpts for the default Opts applied to new instances.
func New(opts ...Opt) *TL {
	tl := &TL{}
	for _, opt := range DefaultOpts {
		opt(tl)
	}
	for _, opt := range opts {
		opt(tl)
	}
	return tl
}

func (tl *TL) Add(t time.Time, entries ...interface{}) {
	for _, entry := range entries {
		tl.items.Insert(valueEntry{entry, t, tl.key})
	}
}

func (tl *TL) Remove(t time.Time, entries ...interface{}) {
	for _, entry := range entries {
		tl.items.Delete(valueEntry{entry, t, tl.key})
	}
}

func (tl *TL) At(t time.Time) []interface{} {
	key := tl.key(t)
	found := tl.items.FindAll(keyEntry(key))
	result := make([]interface{}, len(found))
	for i, element := range found {
		result[i] = entryValue(element)
	}
	return result
}

func (tl *TL) Previous(t time.Time) (previous time.Time, exists bool) {
	key := tl.key(t)
	onOrAfter, ok := tl.items.FindGreaterOrEqual(keyEntry(key))
	last := tl.items.GetLargestNode()
	var prevElement *skiplist.SkipListElement
	if ok {
		prevElement = tl.items.Prev(onOrAfter)
		if prevElement == last {
			// looped around, no entries
			return previous, false
		}
	} else {
		// nothing on or after t, use the largest node
		prevElement = last
	}

	if prevElement == nil {
		// the tl is empty
		return previous, false
	}

	return entryTime(prevElement), true
}

func (tl *TL) Next(t time.Time) (next time.Time, exists bool) {
	key := tl.key(t.Add(1)) // the next possible time
	onOrAfter, ok := tl.items.FindGreaterOrEqual(keyEntry(key))
	if !ok {
		return next, false
	}
	return entryTime(onOrAfter), true
}

func (tl *TL) Bound() (first, last time.Time, exists bool) {
	if tl.Empty() {
		return first, last, false
	}
	firstNode, lastNode := tl.items.GetSmallestNode(), tl.items.GetLargestNode()
	return entryTime(firstNode), entryTime(lastNode), true
}

func (tl *TL) Empty() bool {
	return tl.items.IsEmpty()
}

func entryValue(entry *skiplist.SkipListElement) interface{} {
	value, ok := entry.GetValue().(valueEntry)
	if !ok {
		panic(fmt.Sprintf("entry in skiplist is not a valueEntry?! %T %v", value, value))
	}
	return value.value
}

func entryTime(entry *skiplist.SkipListElement) time.Time {
	value, ok := entry.GetValue().(valueEntry)
	if !ok {
		panic(fmt.Sprintf("entry in skiplist is not a valueEntry?! %T %v", value, value))
	}
	return value.time
}

type valueEntry struct {
	value interface{}
	time  time.Time
	key   KeyFunc
}

func (e valueEntry) ExtractKey() float64 {
	return e.key(e.time)
}

func (e valueEntry) String() string {
	return fmt.Sprintf("%v", e.value)
}

type keyEntry float64

func (e keyEntry) ExtractKey() float64 {
	return float64(e)
}

func (e keyEntry) String() string {
	return fmt.Sprintf("%v", float64(e))
}

// KeyFunc converts a time.Time into a float64 used as a key in the underlying skip list.
type KeyFunc func(t time.Time) float64
