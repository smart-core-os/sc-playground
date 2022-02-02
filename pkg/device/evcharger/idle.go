package evcharger

import (
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	IdleModeID = "idle" // when not charging, including when plugged in or not

	IdleModeMagnitude      = 6
	PluggedInModeMagnitude = 10
)

func (d *Device) setIdle() error {
	startTime := time.Time{}
	var pbStartTime *timestamppb.Timestamp
	if d.model.Unplugged != nil {
		startTime = d.model.Unplugged.At()
		pbStartTime = timestamppb.New(startTime)
	}
	energyLevel := &traits.EnergyLevel{
		PluggedIn: false,
		Flow:      &traits.EnergyLevel_Idle{Idle: &traits.EnergyLevel_Steady{StartTime: pbStartTime}},
	}
	return d.setIdleMode(energyLevel, IdleModeMagnitude, startTime)
}

func (d *Device) setPluggedIn() error {
	energyLevel := &traits.EnergyLevel{
		PluggedIn: true,
		Flow: &traits.EnergyLevel_Idle{Idle: &traits.EnergyLevel_Steady{
			StartTime: timestamppb.New(d.model.PluggedIn.At()),
		}},
		Quantity: d.model.PluggedIn.Level,
	}
	return d.setIdleMode(energyLevel, PluggedInModeMagnitude, d.model.PluggedIn.At())
}

func (d *Device) setIdleMode(energyLevel *traits.EnergyLevel, magnitude float32, startTime time.Time) error {
	if _, err := d.energyStorage.UpdateEnergyLevel(energyLevel); err != nil {
		return err
	}
	mode := d.idleMode(magnitude)
	if !startTime.IsZero() {
		mode.StartTime = timestamppb.New(startTime)
	}
	demand, _ := d.demandAt(d.model.At(), startTime, mode)
	if err := d.updateDemand(demand); err != nil {
		return err
	}
	if err := d.ensureSingleMode(mode); err != nil {
		return err
	}
	if err := d.electric.SetActiveMode(mode); err != nil {
		return err
	}

	return nil
}

// idleMode returns a new traits.ElectricMode that represents _not charging_, i.e. idle.
func (d *Device) idleMode(magnitude float32) *traits.ElectricMode {
	return &traits.ElectricMode{
		Id:          IdleModeID,
		Title:       "Idle",
		Description: "Not charging",
		Normal:      true,
		Segments: []*traits.ElectricMode_Segment{
			{
				Magnitude: magnitude,
				Shape:     &traits.ElectricMode_Segment_Fixed{Fixed: magnitude},
			},
		},
	}
}
