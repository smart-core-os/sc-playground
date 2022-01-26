package timeline

import "time"

// SliceTL is a TL with a Slice method.
type SliceTL interface {
	TL
	// Slice returns a TL containing only events >=from and <to.
	Slice(from, to time.Time) TL
}

// Slice returns a TL containing only tl events >=from and <to.
//
// If tl implements SliceTL, Slice calls tl.Slice.
// Otherwise Slice creates a new TL that excludes events outside the range.
func Slice(tl TL, from, to time.Time) TL {
	if s, ok := tl.(SliceTL); ok {
		return s.Slice(from, to)
	}
	return &sliceTL{tl, from, to}
}

type sliceTL struct {
	tl       TL
	from, to time.Time
}

func (s *sliceTL) Slice(from, to time.Time) TL {
	// if s is entirely contained in [from, to] then nothing will change
	if s.containedWithin(from, to) {
		return s
	}
	return &sliceTL{s.tl, from, to}
}

func (s *sliceTL) containedWithin(from, to time.Time) bool {
	return !from.After(s.from) && !s.to.After(to)
}

func (s *sliceTL) At(t time.Time) []interface{} {
	if t.Before(s.from) || s.to.Before(t) {
		return nil // t outside range, no events
	}
	return s.tl.At(t)
}

func (s *sliceTL) Previous(t time.Time) (previous time.Time, exists bool) {
	if t.Before(s.from) {
		return previous, false // t < from, outside range
	}
	if s.to.Before(t) {
		t = s.to // exclude events >= s.to
	}
	p, ok := s.tl.Previous(t)
	if !ok {
		return p, ok
	}
	if p.Before(s.from) {
		return previous, false // returned event outside our range
	}
	return p, ok
}

func (s *sliceTL) Next(t time.Time) (next time.Time, exists bool) {
	if s.to.Before(t) {
		return next, false // t >= to, outside range
	}
	if t.Before(s.from) {
		t = s.from // exclude events < s.from
	}
	n, ok := s.tl.Next(t)
	if !ok {
		return n, ok
	}
	if !s.to.After(n) {
		return next, false // returned event outside our range
	}
	return n, ok
}

func (s *sliceTL) Bound() (first, last time.Time, exists bool) {
	first, exists = s.min()
	if exists {
		last, _ = s.max()
	}
	return first, last, exists
}

func (s *sliceTL) min() (time.Time, bool) {
	// min event time that's >= from
	if es := s.At(s.from); len(es) > 0 {
		return s.from, true
	}
	return s.Next(s.from)
}

func (s *sliceTL) max() (time.Time, bool) {
	// max event time that < to
	return s.Previous(s.to)
}

func (s *sliceTL) Filter(matches MatchFunc) TL {
	return &sliceTL{Filter(s.tl, matches), s.from, s.to}
}
