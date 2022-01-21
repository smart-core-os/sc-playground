package timeline

import (
	"testing"
	"time"
)

type sliceOnly struct{ SliceTL }

func (o sliceOnly) At(t time.Time) []interface{} {
	return nil
}
func (o sliceOnly) Previous(t time.Time) (previous time.Time, exists bool) {
	return
}
func (o sliceOnly) Next(t time.Time) (next time.Time, exists bool) {
	return
}

func TestSlice(t *testing.T) {
	check := func(name string, sliced, want TL) {
		assertTLEqual(t, name, want, sliced)
		assertTLEqual(t, name+".Slice(Min, Max)", want, Slice(sliced, MinTime, MaxTime))
		assertTLEqual(t, name+".Slice(empty)", Zero(), Slice(sliced, time.Unix(0, 0), time.Unix(0, 0)))
	}

	check("sliceOnly", Slice(sliceOnly{testTL}, time.Unix(2000, 0), time.Unix(5000, 0)), testTL[1:6])
	check("tlOnly", Slice(tlOnly{testTL}, time.Unix(2000, 0), time.Unix(5000, 0)), testTL[1:6])
}
