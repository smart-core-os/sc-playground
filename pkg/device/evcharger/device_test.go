package evcharger

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-playground/internal/th"
	"github.com/smart-core-os/sc-playground/pkg/device"
	"github.com/smart-core-os/sc-playground/pkg/device/evcharger/event"
	"github.com/smart-core-os/sc-playground/pkg/timeline"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDevice_Scrub(t *testing.T) {
	t.Run("empty tl", func(t *testing.T) {
		state := subject(t)
		state.Scrub(1000)
		state.WantDemandMagnitude(IdleModeMagnitude)
		state.WantActiveModeOnly(&traits.ElectricMode{
			Id:          IdleModeID,
			Title:       "Idle",
			Description: "Not charging",
			Normal:      true,
			Segments: []*traits.ElectricMode_Segment{
				{Magnitude: IdleModeMagnitude, Shape: &traits.ElectricMode_Segment_Fixed{Fixed: IdleModeMagnitude}},
			},
		})
		state.WantEnergyLevel(&traits.EnergyLevel{
			PluggedIn: false,
			Flow:      &traits.EnergyLevel_Idle{Idle: &traits.EnergyLevel_Steady{}},
		})
	})
	t.Run("unplugged", func(t *testing.T) {
		state := subject(t,
			event.PlugIn{Event: e(0)},
			event.Unplug{Event: e(1000)},
		)
		state.Scrub(1000)
		state.WantDemandMagnitude(IdleModeMagnitude)
		state.WantActiveModeOnly(&traits.ElectricMode{
			Id:          IdleModeID,
			Title:       "Idle",
			Description: "Not charging",
			Normal:      true,
			StartTime:   pbt(1000),
			Segments: []*traits.ElectricMode_Segment{
				{Magnitude: IdleModeMagnitude, Shape: &traits.ElectricMode_Segment_Fixed{Fixed: IdleModeMagnitude}},
			},
		})
		state.WantEnergyLevel(&traits.EnergyLevel{
			PluggedIn: false,
			Flow:      &traits.EnergyLevel_Idle{Idle: &traits.EnergyLevel_Steady{StartTime: pbt(1000)}},
		})
	})
	t.Run("plugged in", func(t *testing.T) {
		state := subject(t,
			event.Unplug{Event: e(100)},
			event.PlugIn{Event: e(500)},
		)
		state.Scrub(1000)
		state.WantDemandMagnitude(PluggedInModeMagnitude)
		state.WantActiveModeOnly(&traits.ElectricMode{
			Id:          IdleModeID,
			Title:       "Plugged in",
			Description: "Waiting for the charge to start",
			Normal:      true,
			StartTime:   pbt(500),
			Segments: []*traits.ElectricMode_Segment{
				{Magnitude: PluggedInModeMagnitude, Shape: &traits.ElectricMode_Segment_Fixed{Fixed: PluggedInModeMagnitude}},
			},
		})
		state.WantEnergyLevel(&traits.EnergyLevel{
			PluggedIn: true,
			Flow:      &traits.EnergyLevel_Idle{Idle: &traits.EnergyLevel_Steady{StartTime: pbt(500)}},
		})
	})
	t.Run("charging", func(t *testing.T) {
		decisionSegment := &traits.ElectricMode_Segment{Magnitude: DecisionMagnitude, Length: pbd(DecisionDelay), Shape: &traits.ElectricMode_Segment_Fixed{Fixed: DecisionMagnitude}}
		completeSegment := &traits.ElectricMode_Segment{Magnitude: PluggedInModeMagnitude, Shape: &traits.ElectricMode_Segment_Fixed{Fixed: PluggedInModeMagnitude}}
		chargeStart := time.Unix(1000, 0)

		state := subject(t,
			event.PlugIn{
				Event: e(500),
				Full:  &traits.EnergyLevel_Quantity{Percentage: 100, DistanceKm: 300, EnergyKwh: 1000},
				Modes: []*event.ChargeMode{
					{Id: "Normal", Segments: []*traits.ElectricMode_Segment{{Magnitude: 50, Length: pbd(10 * time.Second)}}},
					{Id: "Fast", Segments: []*traits.ElectricMode_Segment{{Magnitude: 100, Length: pbd(6 * time.Second)}}},
					{Id: "Slow", Segments: []*traits.ElectricMode_Segment{{Magnitude: 20, Length: pbd(20 * time.Second)}}},
				},
			},
			event.ChargeStart{Event: e(1000)},
		)

		state.Run("charge start", func(state tester) {
			state.ScrubT(chargeStart)
			state.WantDemandMagnitude(DecisionMagnitude)
			state.WantModesOnly(
				&traits.ElectricMode{
					Id:        "Normal",
					Normal:    true,
					StartTime: pbt(1000),
					Segments: []*traits.ElectricMode_Segment{
						decisionSegment,
						{Magnitude: 50, Length: pbd(10 * time.Second), Shape: &traits.ElectricMode_Segment_Fixed{Fixed: 50}},
						completeSegment,
					},
				},
				&traits.ElectricMode{
					Id:        "Fast",
					StartTime: pbt(1000),
					Segments: []*traits.ElectricMode_Segment{
						decisionSegment,
						{Magnitude: 100, Length: pbd(6 * time.Second), Shape: &traits.ElectricMode_Segment_Fixed{Fixed: 100}},
						completeSegment,
					},
				},
				&traits.ElectricMode{
					Id:        "Slow",
					StartTime: pbt(1000),
					Segments: []*traits.ElectricMode_Segment{
						decisionSegment,
						{Magnitude: 20, Length: pbd(20 * time.Second), Shape: &traits.ElectricMode_Segment_Fixed{Fixed: 20}},
						completeSegment,
					},
				},
			)
			state.WantActiveMode(&traits.ElectricMode{
				Id:        "Normal",
				Normal:    true,
				StartTime: pbt(1000),
				Segments: []*traits.ElectricMode_Segment{
					decisionSegment,
					{Magnitude: 50, Length: pbd(10 * time.Second), Shape: &traits.ElectricMode_Segment_Fixed{Fixed: 50}},
					completeSegment,
				},
			})
			state.WantEnergyLevel(&traits.EnergyLevel{
				PluggedIn: true,
				Flow: &traits.EnergyLevel_Charge{Charge: &traits.EnergyLevel_Transfer{
					StartTime: pbt(1000),
					Time:      pbd(DecisionDelay + 10*time.Second),
					Speed:     traits.EnergyLevel_Transfer_NORMAL,
					Target:    &traits.EnergyLevel_Quantity{Percentage: 100, DistanceKm: 300, EnergyKwh: 1000},
				}},
				Quantity: &traits.EnergyLevel_Quantity{
					Percentage:  0,
					EnergyKwh:   0,
					Descriptive: traits.EnergyLevel_Quantity_EMPTY,
				},
			})
		})

		state.Run("wait for decision", func(state tester) {
			state.ScrubT(chargeStart.Add(DecisionDelay / 2))
			state.WantDemandMagnitude(DecisionMagnitude)
			state.WantModesOnly(
				&traits.ElectricMode{
					Id:        "Normal",
					Normal:    true,
					StartTime: pbt(1000),
					Segments: []*traits.ElectricMode_Segment{
						decisionSegment,
						{Magnitude: 50, Length: pbd(10 * time.Second), Shape: &traits.ElectricMode_Segment_Fixed{Fixed: 50}},
						completeSegment,
					},
				},
				&traits.ElectricMode{
					Id:        "Fast",
					StartTime: pbt(1000),
					Segments: []*traits.ElectricMode_Segment{
						decisionSegment,
						{Magnitude: 100, Length: pbd(6 * time.Second), Shape: &traits.ElectricMode_Segment_Fixed{Fixed: 100}},
						completeSegment,
					},
				},
				&traits.ElectricMode{
					Id:        "Slow",
					StartTime: pbt(1000),
					Segments: []*traits.ElectricMode_Segment{
						decisionSegment,
						{Magnitude: 20, Length: pbd(20 * time.Second), Shape: &traits.ElectricMode_Segment_Fixed{Fixed: 20}},
						completeSegment,
					},
				},
			)
			state.WantActiveMode(&traits.ElectricMode{
				Id:        "Normal",
				Normal:    true,
				StartTime: pbt(1000),
				Segments: []*traits.ElectricMode_Segment{
					decisionSegment,
					{Magnitude: 50, Length: pbd(10 * time.Second), Shape: &traits.ElectricMode_Segment_Fixed{Fixed: 50}},
					completeSegment,
				},
			})
			state.WantEnergyLevel(&traits.EnergyLevel{
				PluggedIn: true,
				Flow: &traits.EnergyLevel_Charge{Charge: &traits.EnergyLevel_Transfer{
					StartTime: pbt(1000),
					Time:      pbd(DecisionDelay + 10*time.Second - DecisionDelay/2),
					Speed:     traits.EnergyLevel_Transfer_NORMAL,
					Target:    &traits.EnergyLevel_Quantity{Percentage: 100, DistanceKm: 300, EnergyKwh: 1000},
				}},
				Quantity: &traits.EnergyLevel_Quantity{
					Percentage:  0,
					EnergyKwh:   0,
					Descriptive: traits.EnergyLevel_Quantity_EMPTY,
				},
			})
		})

		state.Run("half charged", func(state tester) {
			state.ScrubT(chargeStart.Add(DecisionDelay + 5*time.Second))
			state.WantDemandMagnitude(50)
			state.WantActiveModeOnly(&traits.ElectricMode{
				Id:        "Normal",
				Normal:    true,
				StartTime: pbt(1000),
				Segments: []*traits.ElectricMode_Segment{
					decisionSegment,
					{Magnitude: 50, Length: pbd(10 * time.Second), Shape: &traits.ElectricMode_Segment_Fixed{Fixed: 50}},
					completeSegment,
				},
			})
			state.WantEnergyLevel(&traits.EnergyLevel{
				PluggedIn: true,
				Flow: &traits.EnergyLevel_Charge{Charge: &traits.EnergyLevel_Transfer{
					StartTime: pbt(1000),
					Time:      pbd(5 * time.Second),
					Speed:     traits.EnergyLevel_Transfer_NORMAL,
					Target:    &traits.EnergyLevel_Quantity{Percentage: 100, DistanceKm: 300, EnergyKwh: 1000},
				}},
				Quantity: &traits.EnergyLevel_Quantity{
					Percentage:  50,
					DistanceKm:  150,
					EnergyKwh:   500,
					Descriptive: traits.EnergyLevel_Quantity_MEDIUM,
				},
			})
		})
	})
}

func e(t int64) device.Event {
	return device.Event{Created: time.Unix(t, 0)}
}

func pbt(t int64) *timestamppb.Timestamp {
	return timestamppb.New(time.Unix(t, 0))
}

func pbd(d time.Duration) *durationpb.Duration {
	return durationpb.New(d)
}

type tester struct {
	*testing.T
	*Device
}

func subject(t *testing.T, tl ...timeline.Ater) tester {
	return tester{t, New("TEST", timeline.FromSlice(tl), WithDedicatedTL())}
}

func (t tester) Run(name string, fn func(subject tester)) {
	t.T.Helper()
	t.T.Run(name, func(tt *testing.T) {
		fn(tester{tt, t.Device})
	})
}

func (t tester) Scrub(num int64) tester {
	t.T.Helper()
	return t.ScrubT(time.Unix(num, 0))
}

func (t tester) ScrubT(tm time.Time) tester {
	t.T.Helper()
	if err := t.Device.Scrub(tm); err != nil {
		t.T.Fatalf("Scrub(%v) %v", tm, err)
	}
	return t
}

func (t tester) ScrubAdd(d time.Duration) tester {
	t.T.Helper()
	return t.ScrubT(t.model.At().Add(d))
}

func (t tester) WantDemand(wantDemand *traits.ElectricDemand) tester {
	t.T.Helper()
	if diff := th.Diff(wantDemand, t.Device.electric.Demand()); diff != "" {
		t.Fatalf("Demand@%v (-want, +got)\n%v", t.Device.model.At().Unix(), diff)
	}
	return t
}

func (t tester) WantDemandMagnitude(current float32) tester {
	t.T.Helper()
	return t.WantDemand(&traits.ElectricDemand{
		Current: t.Device.magnitudeDemand(t.Device.model.At(), current),
		Rating:  t.Device.rating,
		Voltage: &t.Device.voltage,
	})
}

func (t tester) WantActiveMode(wantMode *traits.ElectricMode) tester {
	t.T.Helper()
	if diff := th.Diff(wantMode, t.Device.electric.ActiveMode()); diff != "" {
		t.Fatalf("ActiveMode@%v (-want, +got)\n%v", t.Device.model.At().Unix(), diff)
	}
	return t
}

func (t tester) WantModesOnly(wantModes ...*traits.ElectricMode) tester {
	t.T.Helper()
	if diff := th.Diff(wantModes, t.Device.electric.Modes(), cmpopts.SortSlices(func(a, b *traits.ElectricMode) bool {
		return a.Id < b.Id
	})); diff != "" {
		t.Fatalf("Modes@%v (-want, +got)\n%v", t.Device.model.At().Unix(), diff)
	}
	return t
}

func (t tester) WantActiveModeOnly(wantMode *traits.ElectricMode) tester {
	t.T.Helper()
	t.WantModesOnly(wantMode)
	t.WantActiveMode(wantMode)
	return t
}

func (t tester) WantEnergyLevel(wantLevel *traits.EnergyLevel) tester {
	t.T.Helper()
	gotLevel, err := t.Device.energyStorage.GetEnergyLevel()
	if err != nil {
		t.Fatal(err)
	}
	if diff := th.Diff(wantLevel, gotLevel); diff != "" {
		t.Fatalf("EnergyLevel@%v (-want, +got)\n%v", t.Device.model.At().Unix(), diff)
	}
	return t
}
