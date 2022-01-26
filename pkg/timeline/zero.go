package timeline

import (
	"time"
)

// Zero returns an empty TL.
// The returned TL will be nil.
func Zero() TL {
	return (*zeroTL)(nil)
}

type zeroTL struct{}

func (z *zeroTL) String() string {
	return "[]"
}

func (z *zeroTL) At(_ time.Time) []interface{} {
	return nil
}

func (z *zeroTL) Previous(_ time.Time) (previous time.Time, exists bool) {
	return previous, false
}

func (z *zeroTL) Next(_ time.Time) (next time.Time, exists bool) {
	return next, false
}

func (z *zeroTL) Empty() bool {
	return true
}

func (z *zeroTL) Bound() (first, last time.Time, exists bool) {
	return first, last, false
}

func (z *zeroTL) Slice(_, _ time.Time) TL {
	return z
}

func (z *zeroTL) Filter(_ MatchFunc) TL {
	return z
}
