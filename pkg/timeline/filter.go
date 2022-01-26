package timeline

import "time"

// MatchFunc returns true if the element at time t should be included in the result.
type MatchFunc func(t time.Time, e interface{}) bool

// FilterTL is a TL with a Filter method.
type FilterTL interface {
	Filter(matches MatchFunc) TL
}

// Filter returns a TL that only includes entries where MatchFunc returns true.
//
// If tl is a FilterTL, Filter returns tl.Filter.
// Otherwise Filter returns a new TL that acts as if no entries exist that MatchFunc returns false for.
func Filter(tl TL, matches MatchFunc) TL {
	if t, ok := tl.(FilterTL); ok {
		return t.Filter(matches)
	}
	return filterTL{tl, matches}
}

type filterTL struct {
	tl      TL
	matches MatchFunc
}

func (f filterTL) At(t time.Time) []interface{} {
	in := f.tl.At(t)
	out := make([]interface{}, 0, len(in))
	for _, item := range in {
		if f.matches(t, item) {
			out = append(out, item)
		}
	}
	return out
}

func (f filterTL) Previous(t time.Time) (previous time.Time, exists bool) {
	for {
		previous, exists := f.tl.Previous(t)
		if exists && f.reallyExists(previous) {
			return previous, exists
		}

		if !exists {
			return previous, false
		}

		t = previous
	}
}

func (f filterTL) Next(t time.Time) (next time.Time, exists bool) {
	for {
		next, exists := f.tl.Next(t)
		if exists && f.reallyExists(next) {
			return next, exists
		}

		if !exists {
			return next, exists
		}

		t = next
	}
}

func (f filterTL) Slice(from, to time.Time) TL {
	return filterTL{Slice(f.tl, from, to), f.matches}
}

func (f filterTL) reallyExists(t time.Time) bool {
	return len(f.At(t)) > 0
}
