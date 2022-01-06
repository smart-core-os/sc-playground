package profile

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smart-core-os/sc-playground/internal/util"
)

func TestProfileFromProto(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name   string
		input  []*traits.ElectricMode_Segment
		expect Profile
	}

	cases := []testCase{
		{
			name: "Converts without a FinalValue",
			input: []*traits.ElectricMode_Segment{
				{
					Length:    durationpb.New(3 * time.Hour),
					Magnitude: 4,
				},
				{
					Length:    durationpb.New(5 * time.Hour),
					Magnitude: 6,
				},
			},
			expect: Profile{
				Segments: []Segment{
					{3 * time.Hour, 4},
					{5 * time.Hour, 6},
				},
				FinalLevel: 0,
			},
		},
		{
			name: "Converts with a FinalValue",
			input: []*traits.ElectricMode_Segment{
				{
					Length:    durationpb.New(3 * time.Hour),
					Magnitude: 4,
				},
				{
					Length:    durationpb.New(5 * time.Hour),
					Magnitude: 6,
				},
				{
					Magnitude: 7,
				},
			},
			expect: Profile{
				Segments: []Segment{
					{3 * time.Hour, 4},
					{5 * time.Hour, 6},
				},
				FinalLevel: 7,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := FromProto(c.input)

			if diff := cmp.Diff(c.expect, actual); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestProfile_Normalised(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name   string
		input  Profile
		expect Profile
	}

	options := cmp.Options{
		cmpopts.EquateEmpty(),
	}

	testCases := []testCase{
		{
			name: "Merging segments of equal duration",
			input: Profile{Segments: []Segment{
				{time.Hour, 1},
				{time.Hour, 1},
			}},
			expect: Profile{Segments: []Segment{
				{2 * time.Hour, 1},
			}},
		},
		{
			name: "Merges segments of unequal duration",
			input: Profile{Segments: []Segment{
				{1 * time.Hour, 1},
				{2 * time.Hour, 1},
			}},
			expect: Profile{Segments: []Segment{
				{3 * time.Hour, 1},
			}},
		},
		{
			name: "Does not merge segments of different levels",
			input: Profile{Segments: []Segment{
				{1 * time.Hour, 1},
				{2 * time.Hour, 2},
			}},
			expect: Profile{Segments: []Segment{
				{1 * time.Hour, 1},
				{2 * time.Hour, 2},
			}},
		},
		{
			name:   "Preserves FinalLevel",
			input:  Profile{FinalLevel: 123},
			expect: Profile{FinalLevel: 123},
		},
		{
			name: "Trailing FinalLevel segments",
			input: Profile{
				Segments: []Segment{
					{1 * time.Hour, 1},
					{2 * time.Hour, 123},
				},
				FinalLevel: 123,
			},
			expect: Profile{
				Segments: []Segment{
					{1 * time.Hour, 1},
				},
				FinalLevel: 123,
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.input.Normalised()
			if diff := cmp.Diff(c.expect, actual, options); diff != "" {
				t.Error(diff)
			}
		})
	}
}

// TestProfile_Normalised_Random checks that Normalised does not affect TotalDuration
func TestProfile_Normalised_Random(t *testing.T) {
	t.Parallel()
	const n = 100

	r := rand.New(util.TestRandSource(t))

	for i := 0; i < n; i++ {
		profile := genProfile(r)

		expected := profile.TotalDuration()
		actual := profile.Normalised().TotalDuration()
		if actual != expected {
			t.Errorf("expected TotalDuration()=%v, but got %v", expected, actual)
		}
	}
}

func TestProfile_MaxLevel(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name   string
		input  Profile
		expect float32
	}

	cases := []testCase{
		{
			name: "Basic operation",
			input: Profile{
				Segments: []Segment{
					{time.Hour, 1},
					{2 * time.Hour, 2},
					{3 * time.Hour, 3},
				},
				FinalLevel: 0,
			},
			expect: 3,
		},
		{
			name: "Empty profile returns FinalLevel",
			input: Profile{
				Segments:   nil,
				FinalLevel: 100,
			},
			expect: 100,
		},
		{
			name: "Zero-duration segments don't count",
			input: Profile{
				Segments: []Segment{
					{0, 123},
					{0, 456},
				},
				FinalLevel: 0,
			},
			expect: 0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.input.MaxLevel()

			if actual != c.expect {
				t.Errorf("expected MaxLevel()=%v, but got %v", c.expect, actual)
			}
		})
	}
}

func TestProfile_LevelAfter(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		profile     Profile
		duration    time.Duration
		expect      float32
		expectPanic bool
	}

	options := cmp.Options{
		cmpopts.EquateEmpty(),
	}

	testCases := []testCase{
		{
			name: "Middle of first segment",
			profile: Profile{Segments: []Segment{
				{2 * time.Hour, 123},
			}},
			duration: time.Hour,
			expect:   123,
		},
		{
			// right on the boundary between two segments, we expect the *second* one to apply
			name: "Start of second segment",
			profile: Profile{Segments: []Segment{
				{1 * time.Hour, 123},
				{1 * time.Hour, 456},
			}},
			duration: time.Hour,
			expect:   456,
		},
		{
			name:        "Rejects negative duration",
			profile:     Profile{},
			duration:    -time.Hour,
			expectPanic: true,
		},
		{
			name: "Past the end",
			profile: Profile{
				Segments: []Segment{
					{time.Hour, 1},
					{time.Hour, 2},
				},
				FinalLevel: 123,
			},
			duration: 3 * time.Hour,
			expect:   123,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			if c.expectPanic {
				defer func() {
					err := recover()
					t.Logf("panic value: %v", err)
				}()
			}

			actual := c.profile.LevelAfter(c.duration)

			if c.expectPanic {
				t.Fatalf("expected test to panic, but it did not; actual = %#v", actual)
			}

			if diff := cmp.Diff(c.expect, actual, options); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestProfile_SplitAt(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		input        Profile
		d            time.Duration
		expectBefore Profile
		expectAfter  Profile
	}

	options := cmp.Options{
		cmpopts.EquateEmpty(),
	}

	cases := []testCase{
		{
			name: "Single segment",
			input: Profile{
				Segments: []Segment{
					{4 * time.Hour, 1},
				},
				FinalLevel: 10,
			},
			d: 1 * time.Hour,
			expectBefore: Profile{
				Segments: []Segment{
					{1 * time.Hour, 1},
				},
				FinalLevel: 0,
			},
			expectAfter: Profile{
				Segments: []Segment{
					{3 * time.Hour, 1},
				},
				FinalLevel: 10,
			},
		},
		{
			name: "On a boundary",
			input: Profile{Segments: []Segment{
				{1 * time.Hour, 1},
				{1 * time.Hour, 2},
			}},
			d: 1 * time.Hour,
			expectBefore: Profile{Segments: []Segment{
				{1 * time.Hour, 1},
			}},
			expectAfter: Profile{Segments: []Segment{
				{1 * time.Hour, 2},
			}},
		},
		{
			name: "Start",
			input: Profile{Segments: []Segment{
				{1 * time.Hour, 1},
			}},
			d:            0,
			expectBefore: Profile{},
			expectAfter: Profile{Segments: []Segment{
				{1 * time.Hour, 1},
			}},
		},
		{
			name: "Zero-length",
			input: Profile{
				Segments: []Segment{
					{0, 1},
				},
				FinalLevel: 2,
			},
			d:            0,
			expectBefore: Profile{},
			expectAfter:  Profile{FinalLevel: 2},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualBefore, actualAfter := c.input.SplitAt(c.d)
			if diff := cmp.Diff(c.expectBefore, actualBefore, options); diff != "" {
				t.Error("before:", diff)
			}
			if diff := cmp.Diff(c.expectAfter, actualAfter, options); diff != "" {
				t.Error("after:", diff)
			}
		})
	}
}

// TestProfile_SplitAt_TotalDuration tests that the total duration of a profile is the same before and after splitting it.
func TestProfile_SplitAt_TotalDuration(t *testing.T) {
	t.Parallel()
	const n = 100

	r := rand.New(util.TestRandSource(t))

	for i := 0; i < n; i++ {
		profile := genProfile(r)
		splitPoint := randDuration(r, profile.TotalDuration())
		a, b := profile.SplitAt(splitPoint)

		expect := profile.TotalDuration()
		actual := a.TotalDuration() + b.TotalDuration()
		if actual != expect {
			t.Errorf("expected total duration to be %v, but it was %v", expect, actual)
		}
	}
}

func TestReduceProfile(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		inputs    []Profile
		expectMax Profile // expected value for Max
		expectSum Profile // expected value for Sum
	}

	cases := []testCase{
		{
			name: "Aligned",
			inputs: []Profile{
				{Segments: []Segment{
					{time.Hour, 1},
				}},
				{Segments: []Segment{
					{time.Hour, 2},
				}},
			},
			expectMax: Profile{
				Segments: []Segment{
					{time.Hour, 2},
				},
			},
			expectSum: Profile{
				Segments: []Segment{
					{time.Hour, 3},
				},
			},
		},
		{
			name: "Overlapping",
			inputs: []Profile{
				{
					Segments: []Segment{
						{1 * time.Hour, 1},
						{2 * time.Hour, 2},
					},
					FinalLevel: 3,
				},
				{
					Segments: []Segment{
						{2 * time.Hour, 4},
					},
					FinalLevel: 5,
				},
			},
			expectMax: Profile{
				Segments: []Segment{
					{2 * time.Hour, 4},
					{1 * time.Hour, 5},
				},
				FinalLevel: 5,
			},
			expectSum: Profile{
				Segments: []Segment{
					{1 * time.Hour, 5},
					{1 * time.Hour, 6},
					{1 * time.Hour, 7},
				},
				FinalLevel: 8,
			},
		},
		{
			name: "FinalLevel",
			inputs: []Profile{
				{
					FinalLevel: 10,
				},
				{
					Segments: []Segment{
						{time.Hour, 1},
						{time.Hour, 2},
					},
					FinalLevel: 10,
				},
			},
			expectMax: Profile{
				Segments: []Segment{
					{2 * time.Hour, 10},
				},
				FinalLevel: 10,
			},
			expectSum: Profile{
				Segments: []Segment{
					{time.Hour, 11},
					{time.Hour, 12},
				},
				FinalLevel: 20,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Run("Max", func(t *testing.T) {
				actual := Max(c.inputs...)

				if diff := cmp.Diff(c.expectMax, actual, cmpopts.EquateEmpty()); diff != "" {
					t.Error(diff)
				}
			})

			t.Run("Sum", func(t *testing.T) {
				actual := Sum(c.inputs...)

				if diff := cmp.Diff(c.expectSum, actual, cmpopts.EquateEmpty()); diff != "" {
					t.Error(diff)
				}
			})
		})
	}
}

func TestSum_Regression(t *testing.T) {
	// this case exposed a bug where Reduce would panic if a profile had any zero-length segments.
	prof := []Profile{
		{
			Segments: []Segment{
				{Duration: 11 * time.Second, Level: 0},
				{Duration: 3 * time.Second, Level: 24},
				{Duration: 10 * time.Second, Level: 60},
			},
			FinalLevel: 48,
		},
		{
			Segments: []Segment{
				{Duration: 7 * time.Second, Level: 0},
			},
			FinalLevel: 48,
		},
		{
			Segments: []Segment{
				{Duration: 0, Level: 0},
			},
			FinalLevel: 48,
		},
	}

	expect := Profile{
		Segments: []Segment{
			{7 * time.Second, 48},
			{4 * time.Second, 96},
			{3 * time.Second, 120},
			{10 * time.Second, 156},
		},
		FinalLevel: 144,
	}

	result := Sum(prof...)

	if diff := cmp.Diff(expect, result, cmpopts.EquateEmpty()); diff != "" {
		t.Error(diff)
	}
}

// TestMaxProfile_Identity verifies that when Max is provided a single Profile, it returns the profile unchanged.
// This is verified using random profiles.
func TestMaxProfile_Identity(t *testing.T) {
	t.Parallel()
	const n = 100

	r := rand.New(util.TestRandSource(t))

	for i := 0; i < n; i++ {
		profile := genProfile(r)

		expect := profile
		actual := Max(profile)

		if diff := cmp.Diff(expect, actual, cmpopts.EquateEmpty()); diff != "" {
			t.Error(diff)
		}
	}
}

func TestProfile_DelayStart(t *testing.T) {
	t.Parallel()
	const n = 100

	r := rand.New(util.TestRandSource(t))

	for i := 0; i < n; i++ {
		profile := genProfile(r)
		delayBy := randDuration(r, time.Hour)

		delayed := profile.DelayStart(delayBy)

		// total duration increases by delayBy
		if expect, actual := profile.TotalDuration()+delayBy, delayed.TotalDuration(); actual != expect {
			t.Errorf("Expected delayed.TotalDuration()=%v, but got %v", expect, actual)
		}

		// number of segments increases by 1
		if expect, actual := len(profile.Segments)+1, len(delayed.Segments); actual != expect {
			t.Errorf("Expected delayed.Segments to have %d elements, but it has %d", expect, actual)
		}

		// first segment has duration == delayBy
		if expect, actual := delayBy, delayed.Segments[0].Duration; actual != expect {
			t.Errorf("Expected first segment to have duration %v, but it was actually %v", expect, actual)
		}

		// first segment has level == 0
		if actual := delayed.Segments[0].Level; actual != 0 {
			t.Errorf("Expected first segment to have Level==0, but it was actually %v", actual)
		}
	}
}

func TestAlign_Random(t *testing.T) {
	r := rand.New(util.TestRandSource(t))

	profiles := make([]Profile, 10)
	for i := range profiles {
		profiles[i] = genProfile(r)
	}

	aligned := Align(profiles)

	t.Run("same time series", func(t *testing.T) {
		for i := range profiles {
			expect := profiles[i].Normalised()
			actual := aligned[i].Normalised()

			if diff := cmp.Diff(expect, actual, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("profile %d: %s", i, diff)
			}
		}
	})

	t.Run("number of segments match", func(t *testing.T) {
		expect := len(aligned[0].Segments)
		for i, p := range aligned {
			actual := len(p.Segments)
			if actual != expect {
				t.Errorf("expected profile %d to have %d segments, but it has %d", i, expect, actual)
			}
		}
	})

	t.Run("segment durations match", func(t *testing.T) {
		for i := range aligned[0].Segments {
			expect := aligned[0].Segments[i].Duration
			for j, p := range aligned {
				actual := p.Segments[i].Duration
				if actual != expect {
					t.Errorf("expected segment %d in profile %d to have duration %v, but it is %v",
						i, j, expect, actual)
				}
			}
		}
	})
}

func genProfile(r *rand.Rand) Profile {
	var profile Profile
	profile.Segments = make([]Segment, r.Intn(20))
	for i := range profile.Segments {
		var segment Segment
		segment.Level = r.Float32() * 100
		segment.Duration = randDuration(r, time.Hour)

		profile.Segments[i] = segment
	}

	return profile
}

func randDuration(r *rand.Rand, max time.Duration) time.Duration {
	if max == 0 {
		return 0
	}
	return time.Duration(r.Int63n(int64(max) + 1))
}
