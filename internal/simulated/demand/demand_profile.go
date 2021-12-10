package demand

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-playground/internal/util/accumulator"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

// Profile represents a piecewise-constant time series of levels. A Profile is a sequence of Segments.
// After the final Segment, the time series has the value FinalLevel for the rest of time.
// Segments that have a duration of 0 are not considered to be part of the time series. They may be removed by any
// operation.
type Profile struct {
	Segments   []Segment
	FinalLevel float32
}

// Segment is a portion of a time series. It represents a constant Level for a given Duration.
type Segment struct {
	Duration time.Duration
	Level    float32
}

func ProfileFromProto(electricMode *traits.ElectricMode) Profile {
	var result Profile

	for _, segment := range electricMode.Segments {
		duration := segment.Length.AsDuration()
		level := segment.Magnitude

		if duration != 0 {
			result.Segments = append(result.Segments, Segment{
				Duration: duration,
				Level:    level,
			})
		}
	}

	if len(electricMode.Segments) != 0 {
		lastSegment := electricMode.Segments[len(electricMode.Segments)-1]
		if lastSegment.Length.AsDuration() == 0 {
			result.FinalLevel = lastSegment.Magnitude
		}
	}

	return result
}

func (p Profile) ToProto() []*traits.ElectricMode_Segment {
	var result []*traits.ElectricMode_Segment

	for _, segment := range p.Segments {
		result = append(result, &traits.ElectricMode_Segment{
			Length:    durationpb.New(segment.Duration),
			Magnitude: segment.Level,
		})
	}

	if p.FinalLevel != 0 {
		result = append(result, &traits.ElectricMode_Segment{
			Magnitude: p.FinalLevel,
		})
	}

	return result
}

// Normalised creates a normalised copy of profile p.
// A normal profile has no zero-duration Segments, and it also does not have adjacent Segments with the same Level.
// The normalised Profile will always return the same results from MaxLevel, TotalDuration and LevelAfter as the
// original profile p.
func (p Profile) Normalised() Profile {
	var (
		oldSegs = p.Segments
		newSegs []Segment
	)

	for len(oldSegs) > 0 {
		var segment Segment
		segment, oldSegs = groupEqualLevels(oldSegs)

		if segment.Duration > 0 {
			newSegs = append(newSegs, segment)
		}
	}

	return Profile{
		Segments:   newSegs,
		FinalLevel: p.FinalLevel,
	}
}

// MaxLevel returns the maximum Level that occurs in any segment
func (p Profile) MaxLevel() float32 {
	max := p.FinalLevel
	for _, segment := range p.Segments {
		if segment.Level > max && segment.Duration > 0 {
			max = segment.Level
		}
	}
	return max
}

func (p Profile) TotalDuration() time.Duration {
	var duration time.Duration = 0
	for _, segment := range p.Segments {
		if segment.Duration < 0 {
			panic("segment has negative duration")
		}
		duration += segment.Duration
	}
	return duration
}

func (p Profile) LevelAfter(d time.Duration) float32 {
	if d < 0 {
		panic("cannot get level after a negative duration")
	}

	// walk the remaining slice until d is less than the duration of the next segment
	remaining := p.Segments
	for len(remaining) > 0 && d >= remaining[0].Duration {
		d -= remaining[0].Duration
		remaining = remaining[1:]
	}

	// if we are within a segment, then use that segment's value, otherwise we are past the end
	if len(remaining) > 0 {
		return remaining[0].Level
	} else {
		return p.FinalLevel
	}
}

func (p Profile) IsEmpty() bool {
	return len(p.Segments) == 0 || p.TotalDuration() == 0
}

func (p Profile) SplitAt(d time.Duration) (before, after Profile) {
	before.FinalLevel = 0
	after.FinalLevel = p.FinalLevel

	segs := p.Segments

	// process whole segments
	for len(segs) > 0 && segs[0].Duration <= d {
		d -= segs[0].Duration
		before.Segments = append(before.Segments, segs[0])
		segs = segs[1:]
	}

	// if the split point is in the middle of a segment, split it in two
	if d > 0 && len(segs) > 0 {
		// partially reduce length of segment
		partial := segs[0]

		if partial.Duration <= d {
			panic("logic error: this whole segment should have been consumed by the loop")
		}

		before.Segments = append(before.Segments, Segment{
			Duration: d,
			Level:    partial.Level,
		})
		after.Segments = append(after.Segments, Segment{
			Duration: partial.Duration - d,
			Level:    partial.Level,
		})

		segs = segs[1:]
	}

	// copy remaining whole segments to after
	after.Segments = append(after.Segments, segs...)

	return
}

// DelayStart returns a copy of p that has a new initial segment, with a level of 0 and duration d.
// The effect is to delay the start of the profile by the given duration.
// Panics if d is negative.
func (p Profile) DelayStart(d time.Duration) Profile {
	if d < 0 {
		panic("cannot delay by negative duration")
	}

	result := p

	// new initial segment
	result.Segments = []Segment{
		{
			Duration: d,
			Level:    0,
		},
	}

	result.Segments = append(result.Segments, p.Segments...)

	return result
}

func MaxProfile(profiles ...Profile) Profile {
	max := func(a, b float32) float32 {
		if a > b {
			return a
		} else {
			return b
		}
	}

	return ReduceProfile(max, profiles...)
}

func SumProfile(profiles ...Profile) Profile {
	sum := func(a, b float32) float32 {
		return a + b
	}

	return ReduceProfile(sum, profiles...)
}

func ReduceProfile(reduce func(acc, level float32) float32, profiles ...Profile) Profile {
	if len(profiles) == 0 {
		panic("cannot reduce over an empty profile list")
	}

	var result Profile

	// the overall final level should be the reduction over all the final levels
	result.FinalLevel = profiles[0].FinalLevel
	for _, profile := range profiles[1:] {
		result.FinalLevel = reduce(result.FinalLevel, profile.FinalLevel)
	}

	// After the end of all a profile's segments, we need to reduce its FinalValue into all the subsequent segments
	// we create. finalAcc accumulates the value to use.
	finalAcc := accumulator.Float32{Reduce: reduce}
	for {
		folded, remaining := reduceInitial(profiles, reduce, &finalAcc)
		if folded.Duration == 0 {
			// there wasn't anything to process
			break
		}

		profiles = remaining

		if finalVal, ok := finalAcc.Get(); ok {
			folded.Level = reduce(folded.Level, finalVal)
		}
		result.Segments = append(result.Segments, folded)
	}

	return result.Normalised()
}

// shortestInitialDuration find the smallest Duration value from among the first segments (i.e. Segments[0]) of
// the provided profiles.
// This can be used to find an appropriate location to split up the profiles' segments, so they can be matched up
// 1-to-1 for further processing.
// Sets ok=true if result is valid. If none of the profiles have any segments, ok=false and the result is undefined.
func shortestInitialDuration(profiles []Profile) (result time.Duration, ok bool) {
	for _, profile := range profiles {
		// skip profiles with no segments
		if len(profile.Segments) == 0 {
			continue
		}

		firstSeg := profile.Segments[0]

		if ok {
			if firstSeg.Duration < result {
				result = firstSeg.Duration
			}
		} else {
			// we have no duration to compare to, so just use this segment as the minimum duration
			result = firstSeg.Duration
			ok = true
		}
	}

	return
}

func reduceInitial(profiles []Profile, reduce func(a, x float32) float32, final *accumulator.Float32) (folded Segment, remaining []Profile) {

	d, ok := shortestInitialDuration(profiles)
	if !ok {
		// if there are no segments in any profile, then return profiles unmodified
		remaining = profiles
		return
	}

	acc := accumulator.Float32{Reduce: reduce}

	for _, profile := range profiles {
		profile = profile.Normalised()
		if profile.IsEmpty() {
			// continue without adding this profile to the remaining profiles slice, to ensure it is not
			// processed further
			final.Accumulate(profile.FinalLevel)
			continue
		}

		prefix, suffix := profile.SplitAt(d)
		if len(prefix.Segments) > 1 {
			panic("logic error: found a segment shorter than d")
		}

		acc.Accumulate(prefix.Segments[0].Level)

		remaining = append(remaining, suffix)
	}

	level, ok := acc.Get()
	if !ok {
		panic("no values to reduce over")
	}

	folded = Segment{
		Duration: d,
		Level:    level,
	}
	return
}

// groupEqualLevels will merge together a contiguous block of segments at the start of the slice that all have the same
// level. The duration of the merged segment is the sum of the durations of the merged segments.
// The remaining unmerged segments, which are unchanged from the profile, will be returned. The returned slice may
// have the same backing array as the segments parameter.
func groupEqualLevels(segments []Segment) (Segment, []Segment) {
	if len(segments) == 0 {
		panic("segments is empty")
	}

	level := segments[0].Level
	grouped := Segment{
		Duration: segments[0].Duration,
		Level:    level,
	}
	segments = segments[1:]

	for len(segments) > 0 {
		segment := segments[0]
		if segment.Level != level {
			return grouped, segments
		}

		grouped.Duration += segment.Duration
		segments = segments[1:]
	}

	return grouped, segments
}
