package evcharger

import (
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-playground/pkg/device/evcharger/event"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	DecisionDelay     = 15 * time.Second
	DecisionMagnitude = PluggedInModeMagnitude
)

func (d *Device) setCharging() error {
	chargeStart := d.model.ChargeStared.At()
	awaitingDecision := d.model.At().Sub(chargeStart) < DecisionDelay

	// update the modes
	modes := d.chargeModesToElectricModes(chargeStart, d.model.PluggedIn.Modes)
	var activeMode *traits.ElectricMode
	if len(d.model.ModeChanges) > 0 {
		// [0] is the most recent mode change
		modeId := d.model.ModeChanges[0].Id
		for _, mode := range modes {
			if mode.Id == modeId {
				activeMode = mode
				break
			}
		}
	}
	if activeMode == nil {
		activeMode = modes[0]
	}

	if awaitingDecision {
		if err := d.ensureOnlyModes(modes); err != nil {
			return err
		}
	} else {
		// only announce that we have a single mode, the active one
		if err := d.ensureSingleMode(activeMode); err != nil {
			return err
		}
	}
	if err := d.electric.SetActiveMode(activeMode); err != nil {
		return err
	}

	// set current demand
	demand, charging := d.demandAt(d.model.At(), chargeStart, activeMode)
	if !charging {
		// charging has ended, likely because the ev is full
		demand = d.magnitudeDemand(d.model.At(), PluggedInModeMagnitude)
	}
	if err := d.updateDemand(demand); err != nil {
		return err
	}

	// update the stored energy for the device
	if _, err := d.energyStorage.UpdateEnergyLevel(d.chargingLevel()); err != nil {
		return err
	}

	return nil
}

func (d *Device) chargeModesToElectricModes(start time.Time, ins []*event.ChargeMode) []*traits.ElectricMode {
	outs := make([]*traits.ElectricMode, len(ins))

	decisionSegment := &traits.ElectricMode_Segment{
		Length:    durationpb.New(DecisionDelay),
		Magnitude: DecisionMagnitude,
		Shape:     &traits.ElectricMode_Segment_Fixed{Fixed: DecisionMagnitude},
	}
	pluggedInSegment := &traits.ElectricMode_Segment{
		Magnitude: PluggedInModeMagnitude,
		Shape:     &traits.ElectricMode_Segment_Fixed{Fixed: PluggedInModeMagnitude},
	}

	for i, in := range ins {
		out := &traits.ElectricMode{}
		out.Id = in.Id
		out.Title = in.Title
		out.Description = in.Description
		out.Normal = i == 0 // priority order, means the first option is the normal one
		out.StartTime = timestamppb.New(start)
		// the first segment of the mode is always the decision segment
		out.Segments = []*traits.ElectricMode_Segment{decisionSegment}
		for _, segment := range in.Segments {
			if segment.Shape == nil {
				segment.Shape = &traits.ElectricMode_Segment_Fixed{Fixed: segment.Magnitude}
			}
			out.Segments = append(out.Segments, segment)
		}
		// the last segment is the device going back into PluggedIn idle mode
		out.Segments = append(out.Segments, pluggedInSegment)

		outs[i] = out
	}
	return outs
}

func (d *Device) chargingLevel() *traits.EnergyLevel {
	// Calculate the currentCharge charge level for the plugged in device
	// as if the active modes segments would charge the device to 100%.
	// Effectively the currentCharge capacity of the storage is the sum of all segment[i].Length * segment[i].Magnitude.
	// Progress is calculated based on how far along the segments we are.
	//
	// We assume charging is complete directly before the first infinite segment is encountered

	// todo: take d.model.PluggedIn.Level into account
	expectedCharge, currentCharge := 0.0, 0.0 // charge times don't include the decision time

	mode := d.electric.ActiveMode()
	now, cur := d.model.At(), mode.StartTime.AsTime()
	for i, segment := range mode.Segments {
		if segment.Length == nil {
			break // first infinite segment, charging complete
		}

		if i > 0 {
			// the first segment is the decision segment, ignore it when calculating totals

			area := segment.Length.AsDuration().Seconds() * float64(segment.Magnitude)
			expectedCharge += area
			if !cur.After(now) {
				segProgress := float64(now.Sub(cur)) / float64(segment.Length.AsDuration())
				if segProgress > 1 {
					// means now is somewhere after this segment ends
					segProgress = 1
				}
				currentCharge += segProgress * area
			}

		}

		cur = cur.Add(segment.Length.AsDuration())
	}

	progress := float32(currentCharge / expectedCharge)
	return &traits.EnergyLevel{
		PluggedIn: true,
		Quantity:  d.chargeQuantity(progress),
		Flow:      &traits.EnergyLevel_Charge{Charge: d.chargeFlow(progress)},
	}
}

func (d *Device) chargeQuantity(progress float32) *traits.EnergyLevel_Quantity {
	var (
		fullPercent  float32 = 100
		fullDistance float32 = 0
		fullEnergy   float32 = 0
	)
	if d.model.PluggedIn.Full != nil {
		full := d.model.PluggedIn.Full
		if full.Percentage > 0 {
			fullPercent = full.Percentage
		}
		if full.DistanceKm > 0 {
			fullDistance = full.DistanceKm
		}
		if full.EnergyKwh > 0 {
			fullEnergy = full.EnergyKwh
		}
	}
	return &traits.EnergyLevel_Quantity{
		Percentage:  progress * fullPercent,
		Descriptive: quantityDescription(progress),
		DistanceKm:  progress * fullDistance,
		EnergyKwh:   progress * fullEnergy,
	}
}

func (d *Device) chargeFlow(chargeProgress float32) *traits.EnergyLevel_Transfer {
	chargeStart := d.model.ChargeStared.At()
	return &traits.EnergyLevel_Transfer{
		StartTime: timestamppb.New(chargeStart),
		Time:      durationpb.New(d.remainingChargeTime()),
		Target:    d.model.PluggedIn.Full,
		Speed:     d.chargeTransferSpeed(chargeProgress),
	}
}

func (d *Device) chargeTransferSpeed(_ float32) traits.EnergyLevel_Transfer_Speed {
	// todo: work out the transfer speed
	//   Idea: do it based on the active mode when compared to other modes
	return traits.EnergyLevel_Transfer_NORMAL
}

func (d *Device) remainingChargeTime() time.Duration {
	chargeStart := d.model.ChargeStared.At()
	now := d.model.At()

	var expectedTimeTaken time.Duration
	for _, segment := range d.electric.ActiveMode().Segments {
		if segment.Length == nil {
			break
		}
		expectedTimeTaken += segment.Length.AsDuration()
	}

	takenTime := now.Sub(chargeStart)
	duration := expectedTimeTaken - takenTime
	if duration < 0 {
		return 0
	}
	return duration
}
