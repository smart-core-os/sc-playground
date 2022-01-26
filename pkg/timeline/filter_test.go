package timeline

import (
	"testing"
	"time"
)

type filterOnly struct{ FilterTL }

func (o filterOnly) At(t time.Time) []interface{} {
	return nil
}
func (o filterOnly) Previous(t time.Time) (previous time.Time, exists bool) {
	return
}
func (o filterOnly) Next(t time.Time) (next time.Time, exists bool) {
	return
}

func TestFilter(t *testing.T) {
	check := func(t *testing.T, name string, filtered, want TL) {
		assertTLEqual(t, name, want, filtered)
		assertTLEqual(t, name+".Filter(All)", want, Filter(filtered, matchesAll()))
		assertTLEqual(t, name+".Filter(None)", Zero(), Filter(filtered, matchesNone()))
	}

	filtered := aterTL{
		testTL[0], // 1000
		testTL[3], // 3000
		testTL[4], // 4000
		testTL[5], // 4000
	}
	matches := matchesOnly(
		time.Unix(1000, 0),
		time.Unix(3000, 0),
		time.Unix(4000, 0),
	)

	t.Run("filterOnly", func(t *testing.T) {
		check(t, "filterOnly", Filter(filterOnly{testTL}, matches), filtered)
	})
	t.Run("tlOnly", func(t *testing.T) {
		check(t, "tlOnly", Filter(tlOnly{testTL}, matches), filtered)
	})
}

func matchesOnly(ts ...time.Time) MatchFunc {
	return func(t time.Time, e interface{}) bool {
		for _, candidate := range ts {
			if t.Equal(candidate) {
				return true
			}
		}
		return false
	}
}
