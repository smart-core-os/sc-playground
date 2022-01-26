package timeline

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAterTL_Empty(t *testing.T) {
	tests := []struct {
		name string
		tl   aterTL
		want bool
	}{
		{"empty", aterTL([]Ater{}), true},
		{"one", aterTL([]Ater{atUnix(0, 0)}), false},
		{"testTL", testTL, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tl.Empty(); got != tt.want {
				t.Fatalf("aterTL.Empty() want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestAterTL_Bound(t *testing.T) {
	type want struct {
		from, to time.Time
		exists   bool
	}
	tests := []struct {
		name string
		tl   aterTL
		want want
	}{
		{"empty", aterTL([]Ater{}), want{exists: false}},
		{"one", aterTL([]Ater{atUnix(0, 0)}), want{time.Unix(0, 0), time.Unix(0, 0), true}},
		{"one later", aterTL([]Ater{atUnix(1000, 0)}), want{time.Unix(1000, 0), time.Unix(1000, 0), true}},
		{"testTL", testTL, want{testTL[0].At(), testTL[len(testTL)-1].At(), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if from, to, exists := tt.tl.Bound(); !from.Equal(tt.want.from) || !to.Equal(tt.want.to) || exists != tt.want.exists {
				t.Fatalf("aterTL.Bound() want %v, got %v", tt.want, want{from, to, exists})
			}
		})
	}
}

func TestAterTL_At(t *testing.T) {
	tests := []struct {
		name string
		tl   aterTL
		at   time.Time
		want []interface{}
	}{
		{"empty", aterTL([]Ater{}), time.Unix(0, 0), nil},
		{"empty later", aterTL([]Ater{}), time.Unix(1000, 0), nil},
		{"empty MinTime", aterTL([]Ater{}), MinTime, nil},
		{"empty MaxTime", aterTL([]Ater{}), MaxTime, nil},
		{"one", aterTL([]Ater{atUnix(0, 0)}), time.Unix(0, 0), []interface{}{atUnix(0, 0)}},
		{"one missing", aterTL([]Ater{atUnix(0, 0)}), time.Unix(0, 1), nil},
		{"testTL missing", testTL, time.Unix(1000, 1), nil},
		{"testTL single", testTL, time.Unix(1000, 0), []interface{}{atUnix(1000, 0)}},
		{"testTL multiple", testTL, time.Unix(2000, 0), []interface{}{atUnix(2000, 0), atUnix(2000, 0)}},
		{"testTL MinTime", testTL, MinTime, nil},
		{"testTL MaxTime", testTL, MaxTime, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tt.tl.At(tt.at), cmpTimeAtAsTime()); diff != "" {
				t.Fatalf("aterTL.At(%v) (-want, +got)\n%v", tt.at, diff)
			}
		})
	}
}

func TestAterTL_Previous(t *testing.T) {
	type want struct {
		t      time.Time
		exists bool
	}
	tests := []struct {
		name string
		tl   aterTL
		at   time.Time
		want want
	}{
		{"empty", aterTL([]Ater{}), time.Unix(0, 0), want{exists: false}},
		{"empty later", aterTL([]Ater{}), time.Unix(1000, 0), want{exists: false}},
		{"empty MinTime", aterTL([]Ater{}), MinTime, want{exists: false}},
		{"empty MaxTime", aterTL([]Ater{}), MaxTime, want{exists: false}},
		{"one at", aterTL([]Ater{atUnix(0, 0)}), time.Unix(0, 0), want{exists: false}},
		{"one before", aterTL([]Ater{atUnix(0, 0)}), time.Unix(0, -1), want{exists: false}},
		{"one after", aterTL([]Ater{atUnix(0, 0)}), time.Unix(0, 1), want{time.Unix(0, 0), true}},
		{"testTL before first", testTL, time.Unix(0, 0), want{exists: false}},
		{"testTL at first", testTL, time.Unix(1000, 0), want{exists: false}},
		{"testTL after first", testTL, time.Unix(1000, 1), want{time.Unix(1000, 0), true}},
		{"testTL at second", testTL, time.Unix(2000, 0), want{time.Unix(1000, 0), true}},
		{"testTL prev is multiple", testTL, time.Unix(3000, 0), want{time.Unix(2000, 0), true}},
		{"testTL after end", testTL, time.Unix(10000, 0), want{time.Unix(5000, 0), true}},
		{"testTL MinTime", testTL, MinTime, want{exists: false}},
		{"testTL MaxTime", testTL, MaxTime, want{time.Unix(5000, 0), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if p, exists := tt.tl.Previous(tt.at); exists != tt.want.exists || !p.Equal(tt.want.t) {
				t.Fatalf("aterTL.Previous(%v) want %v, got %v", tt.at, tt.want, want{p, exists})
			}
		})
	}
}

func TestAterTL_Next(t *testing.T) {
	type want struct {
		t      time.Time
		exists bool
	}
	tests := []struct {
		name string
		tl   aterTL
		at   time.Time
		want want
	}{
		{"empty", aterTL([]Ater{}), time.Unix(0, 0), want{exists: false}},
		{"empty later", aterTL([]Ater{}), time.Unix(1000, 0), want{exists: false}},
		{"empty MinTime", aterTL([]Ater{}), MinTime, want{exists: false}},
		{"empty MaxTime", aterTL([]Ater{}), MaxTime, want{exists: false}},
		{"one at", aterTL([]Ater{atUnix(0, 0)}), time.Unix(0, 0), want{exists: false}},
		{"one before", aterTL([]Ater{atUnix(0, 0)}), time.Unix(0, -1), want{time.Unix(0, 0), true}},
		{"one after", aterTL([]Ater{atUnix(0, 0)}), time.Unix(0, 1), want{exists: false}},
		{"testTL before first", testTL, time.Unix(0, 0), want{time.Unix(1000, 0), true}},
		{"testTL after first", testTL, time.Unix(1000, 1), want{time.Unix(2000, 0), true}},
		{"testTL next is multiple", testTL, time.Unix(1000, 0), want{time.Unix(2000, 0), true}},
		{"testTL before end", testTL, time.Unix(5000, -1), want{time.Unix(5000, 0), true}},
		{"testTL at end", testTL, time.Unix(5000, 0), want{exists: false}},
		{"testTL after end", testTL, time.Unix(10000, 0), want{exists: false}},
		{"testTL MinTime", testTL, MinTime, want{time.Unix(1000, 0), true}},
		{"testTL MaxTime", testTL, MaxTime, want{exists: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if p, exists := tt.tl.Next(tt.at); exists != tt.want.exists || !p.Equal(tt.want.t) {
				t.Fatalf("aterTL.Next(%v) want %v, got %v", tt.at, tt.want, want{p, exists})
			}
		})
	}
}

func TestAterTL_Slice(t *testing.T) {
	tests := []struct {
		name string
		tl   aterTL
		from time.Time
		to   time.Time
		want TL
	}{
		{"empty", aterTL([]Ater{}), time.Unix(0, 0), time.Unix(100, 0), Zero()},
		{"one before", aterTL([]Ater{atUnix(1000, 0)}), time.Unix(0, 0), time.Unix(100, 0), Zero()},
		{"one upto", aterTL([]Ater{atUnix(1000, 0)}), time.Unix(0, 0), time.Unix(1000, 0), Zero()},
		{"one after", aterTL([]Ater{atUnix(1000, 0)}), time.Unix(1500, 0), time.Unix(2000, 0), Zero()},
		{"one following", aterTL([]Ater{atUnix(1000, 0)}), time.Unix(1000, 0), time.Unix(2000, 0), FromSlice([]Ater{atUnix(1000, 0)})},
		{"one around", aterTL([]Ater{atUnix(1000, 0)}), time.Unix(500, 0), time.Unix(2000, 0), FromSlice([]Ater{atUnix(1000, 0)})},
		{"testTL empty range", testTL, time.Unix(1000, 0), time.Unix(1000, 0), Zero()},
		{"testTL before", testTL, time.Unix(0, 0), time.Unix(1000, 0), Zero()},
		{"testTL after", testTL, time.Unix(5000, 1), time.Unix(10000, 0), Zero()},
		{"testTL first", testTL, time.Unix(500, 1), time.Unix(1500, 0), FromSlice([]Ater{atUnix(1000, 0)})},
		{"testTL first (not second)", testTL, time.Unix(500, 1), time.Unix(2000, 0), FromSlice([]Ater{atUnix(1000, 0)})},
		{"testTL multiple", testTL, time.Unix(2000, 0), time.Unix(3000, 0), FromSlice([]Ater{atUnix(2000, 0), atUnix(2000, 0)})},
		{"testTL [1,3)", testTL, time.Unix(1000, 0), time.Unix(3000, 0), FromSlice([]Ater{atUnix(1000, 0), atUnix(2000, 0), atUnix(2000, 0)})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := fmt.Sprintf("aterTL.Slice(%v, %v)", tt.from, tt.to)
			assertTLEqual(t, name, tt.want, tt.tl.Slice(tt.from, tt.to))
		})
	}
}

func TestAterTL_Filter(t *testing.T) {
	tests := []struct {
		name    string
		tl      aterTL
		matches MatchFunc
		want    TL
	}{
		{"empty", aterTL([]Ater{}), matchesAll(), Zero()},
		{"one include", aterTL([]Ater{atUnix(1000, 0)}), matchesAll(), FromSlice([]Ater{atUnix(1000, 0)})},
		{"one exclude", aterTL([]Ater{atUnix(1000, 0)}), matchesNone(), Zero()},
		{"at same time", aterTL([]Ater{namedAtUnix("A", 1000, 0), namedAtUnix("B", 1000, 0)}),
			matchesName("B"), FromSlice([]Ater{namedAtUnix("B", 1000, 0)})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTLEqual(t, "aterTL.Filter", tt.want, tt.tl.Filter(tt.matches))
		})
	}
}

func matchesAll() MatchFunc {
	return func(t time.Time, e interface{}) bool {
		return true
	}
}

func matchesNone() MatchFunc {
	return func(t time.Time, e interface{}) bool {
		return false
	}
}

func matchesName(n string) MatchFunc {
	return func(t time.Time, e interface{}) bool {
		if named, ok := e.(interface{ Name() string }); ok {
			return named.Name() == n
		}
		return false
	}
}

func namedAtUnix(n string, sec, nsec int64) namedAt {
	return namedAt{time.Unix(sec, nsec), n}
}

type namedAt struct {
	T time.Time
	N string
}

func (n namedAt) At() time.Time {
	return n.T
}

func (n namedAt) Name() string {
	return n.N
}
