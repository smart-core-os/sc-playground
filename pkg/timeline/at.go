package timeline

import (
	"sort"
	"time"
)

// Ater provides the At method for converting to a Time.
type Ater interface {
	// At returns the time this object is at.
	At() time.Time
}

// FromSlice creates a TL using a []Ater.
// This will sort s in place by the results of calling s[i].At().
func FromSlice(s []Ater) TL {
	sort.Slice(s, func(i, j int) bool {
		return s[i].At().Before(s[j].At())
	})
	return aterTL(s)
}

type aterTL []Ater

func (s aterTL) search(t time.Time) int {
	return sort.Search(len(s), func(i int) bool {
		return !s[i].At().Before(t)
	})
}

func (s aterTL) Empty() bool {
	return len(s) == 0
}

func (s aterTL) At(t time.Time) []interface{} {
	i := s.search(t)
	if i == len(s) {
		return nil
	}

	var res []interface{}
	for ; i < len(s) && s[i].At().Equal(t); i++ {
		res = append(res, s[i])
	}
	return res
}

func (s aterTL) Previous(t time.Time) (previous time.Time, exists bool) {
	i := s.search(t)
	if i == 0 {
		return previous, false
	}
	return s[i-1].At(), true
}

func (s aterTL) Next(t time.Time) (next time.Time, exists bool) {
	_, next, exists = s.next(t)
	return next, exists
}

func (s aterTL) next(t time.Time) (i int, next time.Time, exists bool) {
	i = s.search(t)
	for ; i < len(s); i++ {
		r := s[i].At()
		if !r.Equal(t) {
			return i, r, true
		}
	}
	return i, next, false
}

func (s aterTL) Bound() (first, last time.Time, exists bool) {
	if s.Empty() {
		return first, last, false
	}
	return s[0].At(), s[len(s)-1].At(), true
}

func (s aterTL) Slice(from, to time.Time) TL {
	fi := s.search(from)
	ti := s.search(to)
	return s[fi:ti]
}

// Filter implements FilterTL.
// This implementation eagerly applies matches to all entries in s.
func (s aterTL) Filter(matches MatchFunc) TL {
	var out aterTL
	for _, e := range s {
		if matches(e.At(), e) {
			out = append(out, e)
		}
	}
	return out
}
