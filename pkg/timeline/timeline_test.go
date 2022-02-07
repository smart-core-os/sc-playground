package timeline

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestMinMaxTime(t *testing.T) {
	if !MinTime.Before(MaxTime) {
		t.Fatalf("MinTime not < MaxTime: %v, %v", MinTime, MaxTime)
	}
	if !MaxTime.After(MinTime) {
		t.Fatalf("MaxTime not > MinTime: %v, %v", MaxTime, MinTime)
	}
	if !MinTime.Before(time.Unix(0, 0)) {
		t.Fatalf("MinTime not before Unix(0, 0): %v, %v", MinTime, time.Unix(0, 0))
	}
	if !MaxTime.After(time.Unix(0, 0)) {
		t.Fatalf("MaxTime not after Unix(0, 0): %v, %v", MaxTime, time.Unix(0, 0))
	}
}

var testTL = aterTL{
	timeAt(time.Unix(1000, 0)),
	timeAt(time.Unix(2000, 0)),
	timeAt(time.Unix(2000, 0)),
	timeAt(time.Unix(3000, 0)),
	timeAt(time.Unix(4000, 0)),
	timeAt(time.Unix(4000, 0)),
	timeAt(time.Unix(5000, 0)),
}

type timeAt time.Time

func (t timeAt) At() time.Time {
	return time.Time(t)
}
func (t timeAt) String() string {
	return time.Time(t).String()
}

func atUnix(sec, nsec int64) timeAt {
	return timeAt(time.Unix(sec, nsec))
}

type tlOnly struct{ TL }

func cmpTimeAtAsTime() cmp.Option {
	return cmp.Transformer("time.Time", func(at timeAt) time.Time {
		return time.Time(at)
	})
}

func assertTLEqual(t *testing.T, name string, a, b TL) {
	t.Helper()
	if equal, reason := EqualExplain(a, b); !equal {
		t.Fatalf("%v %v", name, reason)
	}
}
