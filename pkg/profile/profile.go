package profile

import (
	"sort"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Profile represents a piecewise-constant time series of levels. A Profile is a sequence of Segments.
// After the final Segment, the time series has the value FinalLevel for the rest of time.
// Segments that have a duration of 0 are not considered to be part of the time series. They may be removed by any
// operation.
type Profile struct {
	// Segments are the pieces of the profile, stored in time order. They form a time series, with Segments[0] starting
	// at time 0, and other segments starting immediately after the previous segment's duration has elapsed.
	Segments []Segment
	// FinalLevel is the level of this profile after all Segments are completed. After TotalDuration has elapsed,
	// the Profile's level will be FinalLevel forever.
	FinalLevel float32
}

// Segment is a portion of a time series. It represents a constant Level for a given Duration.
type Segment struct {
	// Duration is how long this Segment is in effect for. It's relative to the end of the previous Segment
	// in a Profile. Duration must be non-negative.
	Duration time.Duration
	Level    float32
}

// FromProto constructs a profile from the segments in a Smart Core ElectricMode.
// Levels are taken from ElectricMode_Segment.Magnitude.
// The FinalLevel is the Magnitude of the final segment in the ElectricMode, if that segment has 0 duration.
// Otherwise, FinalLevel will be 0.
func FromProto(electricSegments []*traits.ElectricMode_Segment) Profile {
	var result Profile

	for _, segment := range electricSegments {
		duration := segment.Length.AsDuration()
		level := segment.Magnitude

		if duration != 0 {
			result.Segments = append(result.Segments, Segment{
				Duration: duration,
				Level:    level,
			})
		}
	}

	if len(electricSegments) != 0 {
		lastSegment := electricSegments[len(electricSegments)-1]
		if lastSegment.Length.AsDuration() == 0 {
			result.FinalLevel = lastSegment.Magnitude
		}
	}

	return result
}

// ToProto converts p into a slice of segments that can be used in the Segments field of an ElectricMode.
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
// The Normalised profile represents the same time series as p, but removes redundant segments.
// A normal profile has no zero-duration Segments, and it also does not have adjacent Segments with the same Level.
// Segments at the end of the Profile with the same Level as FinalLevel will also be removed.
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

	// remove trailing segments, if they are equal to FinalLevel
	for i := len(newSegs) - 1; i >= 0; i-- {
		if newSegs[i].Level == p.FinalLevel {
			newSegs = newSegs[:i]
		} else {
			break
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

// TotalDuration returns the sum of the Duration in all Segments.
// Panics if any Segment has a negative Duration.
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

// LevelAfter returns the Level of the segment that is in effect at duration d past the start time of the Profile.
// If d >= TotalDuration, then FinalLevel is returned.
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

// IsEmpty returns true if there are no segments or all segments have 0 Duration, false otherwise.
func (p Profile) IsEmpty() bool {
	return len(p.Segments) == 0 || p.TotalDuration() == 0
}

// Truncate removes part of Profile p, from the start, totalling duration d.
// If d is greater than or equal to the TotalDuration of p, then a zero-duration Profile with the same FinalLevel as
// p is returned.
// If d is zero, then the Profile returned is equal to p.
// d must be non-negative.
func (p Profile) Truncate(d time.Duration) Profile {
	if d >= p.TotalDuration() {
		return Profile{
			Segments:   nil,
			FinalLevel: p.FinalLevel,
		}
	}

	result := Profile{FinalLevel: p.FinalLevel}
	skip := d // amount of time left to discard
	for _, segment := range p.Segments {
		if skip >= segment.Duration {
			skip -= segment.Duration
			continue // skip this segment entirely
		}

		switch {
		case skip == 0:
			// keep the entire segment
			result.Segments = append(result.Segments, segment)
		case skip > 0:
			// shorten the segment appropriately, and add it
			keep := segment.Duration - skip
			result.Segments = append(result.Segments, Segment{
				Duration: keep,
				Level:    segment.Level,
			})

			skip = 0 // this segment is the first included one; don't need to skip any more
		default:
			panic("logic error: invalid skip value")
		}
	}

	return result
}

// SplitAt splits a Profile into two parts at duration d after the start of the profile.
// The before profile contains all segments entirely before d, and the after profile contains all segments entirely after
// d. The segment on the boundary is split between before and after.
// before.TotalDuration() == d
func (p Profile) SplitAt(d time.Duration) (before, after Profile) {
	before.FinalLevel = 0
	after.FinalLevel = p.FinalLevel

	total := p.TotalDuration()
	if d > total {
		extra := d - total
		// Include the extra beyond the end in before,
		// so that before.TotalDuration() == d even if p is shorter than that
		// This will mean after is empty
		before.Segments = make([]Segment, len(p.Segments))
		copy(before.Segments, p.Segments)

		before.Segments = append(before.Segments, Segment{
			Duration: extra,
			Level:    p.FinalLevel,
		})
		return
	}

	segs := p.Segments

	// process whole segments
	for len(segs) > 0 && segs[0].Duration <= d {
		d -= segs[0].Duration
		if segs[0].Duration > 0 {
			before.Segments = append(before.Segments, segs[0])
		}
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
	result.Segments = make([]Segment, len(p.Segments)+1)
	copy(result.Segments[1:], p.Segments)

	// new initial segment
	result.Segments[0] = Segment{
		Duration: d,
		Level:    0,
	}

	return result
}

// Max reduces Profiles by taking the maximum Level at every instant.
func Max(profiles ...Profile) Profile {
	max := func(a, b float32) float32 {
		if a > b {
			return a
		} else {
			return b
		}
	}

	return Reduce(max, profiles...)
}

// Sum adds profiles together, such that the level at every time is the sum of the levels of each profile at that time.
func Sum(profiles ...Profile) Profile {
	sum := func(a, b float32) float32 {
		return a + b
	}

	return Reduce(sum, profiles...)
}

// Reduce will combine profiles together using a provided reduction function.
// The reduction function (which should be deterministic) is used to combine the values from multiple profiles together.
// In the returned profile, the level at any instant is the reduction of all the Profiles' levels at that same instant.
func Reduce(reduce func(acc, level float32) float32, profiles ...Profile) Profile {
	if len(profiles) == 0 {
		panic("cannot reduce over an empty profile list")
	}

	var result Profile

	// the overall final level should be the reduction over all the final levels
	result.FinalLevel = profiles[0].FinalLevel
	for _, profile := range profiles[1:] {
		result.FinalLevel = reduce(result.FinalLevel, profile.FinalLevel)
	}

	// align the segment starts and ends of all profiles
	profiles = Align(profiles)
	// all profiles should now have the same number of segments, and of identical duration
	result.Segments = make([]Segment, len(profiles[0].Segments))

	for i := range result.Segments {
		reduced := Segment{
			Duration: profiles[0].Segments[i].Duration,
			Level:    profiles[0].Segments[i].Level,
		}

		// reduce the values of the aligned segments in the current position
		for _, p := range profiles[1:] {
			segment := p.Segments[i]
			reduced.Level = reduce(reduced.Level, segment.Level)
		}

		result.Segments = append(result.Segments, reduced)
	}

	return result.Normalised()
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

// Align will produce a slice of Profiles representing the same time series as the input, but with all profiles
// having the same number of segments and all segments in the same position having the same duration.
// This makes it easier to
func Align(profiles []Profile) []Profile {
	edges := segmentEdges(profiles)
	results := make([]Profile, 0, len(profiles))

	for _, p := range profiles {
		result := p
		result.Segments = make([]Segment, 0, len(edges))

		for j, start := range edges[:len(edges)-1] {
			end := edges[j+1]
			length := end - start
			segment := Segment{
				Duration: length,
				Level:    p.LevelAfter(start),
			}
			result.Segments = append(result.Segments, segment)
		}

		results = append(results, result)
	}

	return results
}

func segmentEdges(profiles []Profile) []time.Duration {
	// gather all of the segment start times, as offsets from the start of their respective profiles
	var edges []time.Duration
	for _, p := range profiles {
		var offset time.Duration = 0
		for _, s := range p.Segments {
			edges = append(edges, offset)
			offset += s.Duration
		}

		// also add the end of the last segment
		edges = append(edges, offset)
	}

	// remove duplicates
	seen := make(map[time.Duration]bool)
	i := 0
	for _, start := range edges {
		_, present := seen[start]
		if !present {
			edges[i] = start
			seen[start] = true
			i++
		}
	}
	edges = edges[:i]

	// sort the start times
	sort.Slice(edges, func(i, j int) bool {
		return edges[i] < edges[j]
	})

	return edges
}
