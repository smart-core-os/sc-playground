package energystorage

import (
	"testing"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
	"github.com/smart-core-os/sc-playground/internal/th"
	"github.com/tanema/gween/ease"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestOscillator_Scrub(t *testing.T) {
	// lets make a nice easy oscillator we can test against
	start := time.UnixMilli(int64(12 * time.Hour))
	model := energystorage.NewModel()
	o := NewOscillator(
		model,
		WithCycleStart(start),
		WithChargeDuration(2*time.Minute),
		WithFullDuration(4*time.Minute),
		WithDischargeDuration(6*time.Minute),
		WithEmptyDuration(8*time.Minute),
		WithChargeRamp(ease.Linear),
		WithDischargeRamp(ease.Linear),
	)
	tests := []struct {
		name string
		t    time.Time
		want *traits.EnergyLevel
	}{
		{name: "charging 0%", t: start, want: &traits.EnergyLevel{
			PluggedIn: true,
			Quantity:  &traits.EnergyLevel_Quantity{},
			Flow: &traits.EnergyLevel_Charge{Charge: &traits.EnergyLevel_Transfer{
				Time:      durationpb.New(2 * time.Minute),
				StartTime: timestamppb.New(start),
				Target: &traits.EnergyLevel_Quantity{
					Percentage:  100,
					Descriptive: traits.EnergyLevel_Quantity_FULL,
				},
			}},
		}},
		{name: "charging 50%", t: start.Add(time.Minute), want: &traits.EnergyLevel{
			PluggedIn: true,
			Quantity:  &traits.EnergyLevel_Quantity{Percentage: 50},
			Flow: &traits.EnergyLevel_Charge{Charge: &traits.EnergyLevel_Transfer{
				Time:      durationpb.New(time.Minute),
				StartTime: timestamppb.New(start),
				Target: &traits.EnergyLevel_Quantity{
					Percentage:  100,
					Descriptive: traits.EnergyLevel_Quantity_FULL,
				},
			}},
		}},
		{name: "charging 75%", t: start.Add(90 * time.Second), want: &traits.EnergyLevel{
			PluggedIn: true,
			Quantity:  &traits.EnergyLevel_Quantity{Percentage: 75},
			Flow: &traits.EnergyLevel_Charge{Charge: &traits.EnergyLevel_Transfer{
				Time:      durationpb.New(30 * time.Second),
				StartTime: timestamppb.New(start),
				Target: &traits.EnergyLevel_Quantity{
					Percentage:  100,
					Descriptive: traits.EnergyLevel_Quantity_FULL,
				},
			}},
		}},
		{name: "full", t: start.Add(2 * time.Minute), want: &traits.EnergyLevel{
			PluggedIn: true,
			Quantity: &traits.EnergyLevel_Quantity{
				Percentage:  100,
				Descriptive: traits.EnergyLevel_Quantity_FULL,
			},
			Flow: &traits.EnergyLevel_Idle{Idle: &traits.EnergyLevel_Steady{
				StartTime: timestamppb.New(start.Add(2 * time.Minute)),
			}},
		}},
		{name: "discharge 0%", t: start.Add(6 * time.Minute), want: &traits.EnergyLevel{
			PluggedIn: true,
			Quantity:  &traits.EnergyLevel_Quantity{Percentage: 100},
			Flow: &traits.EnergyLevel_Discharge{Discharge: &traits.EnergyLevel_Transfer{
				Time:      durationpb.New(6 * time.Minute),
				StartTime: timestamppb.New(start.Add(6 * time.Minute)),
			}},
		}},
		{name: "discharge 33%", t: start.Add(8 * time.Minute), want: &traits.EnergyLevel{
			PluggedIn: true,
			Quantity:  &traits.EnergyLevel_Quantity{Percentage: 66.666666},
			Flow: &traits.EnergyLevel_Discharge{Discharge: &traits.EnergyLevel_Transfer{
				Time:      durationpb.New(4 * time.Minute),
				StartTime: timestamppb.New(start.Add(6 * time.Minute)),
			}},
		}},
		{name: "empty", t: start.Add(12 * time.Minute), want: &traits.EnergyLevel{
			PluggedIn: true,
			Quantity: &traits.EnergyLevel_Quantity{
				Descriptive: traits.EnergyLevel_Quantity_EMPTY,
			},
			Flow: &traits.EnergyLevel_Idle{Idle: &traits.EnergyLevel_Steady{
				StartTime: timestamppb.New(start.Add(12 * time.Minute)),
			}},
		}},
		{name: "charging 50% again", t: start.Add(21 * time.Minute), want: &traits.EnergyLevel{
			PluggedIn: true,
			Quantity:  &traits.EnergyLevel_Quantity{Percentage: 50},
			Flow: &traits.EnergyLevel_Charge{Charge: &traits.EnergyLevel_Transfer{
				Time:      durationpb.New(time.Minute),
				StartTime: timestamppb.New(start.Add(20 * time.Minute)),
				Target: &traits.EnergyLevel_Quantity{
					Percentage:  100,
					Descriptive: traits.EnergyLevel_Quantity_FULL,
				},
			}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := o.Scrub(tt.t); err != nil {
				t.Errorf("Scrub() error = %v", err)
			}
			got, _ := model.GetEnergyLevel()
			if diff := th.Diff(tt.want, got); diff != "" {
				t.Errorf("Scrub() (-want,+got) diff = %v", diff)
			}
		})
	}
}
