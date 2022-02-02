package evcharger

import (
	"fmt"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
)

// ensureSingleMode adjusts the electric modes of the device such that only the given mode exists.
func (d *Device) ensureSingleMode(mode *traits.ElectricMode) error {
	modes := d.electric.Modes(nil)
	modeExists := false
	for _, m := range modes {
		if m.Id == mode.Id {
			modeExists = true
			if _, err := d.electric.UpdateMode(mode, nil); err != nil {
				return err
			}
		} else {
			if err := d.electric.DeleteMode(m.Id); err != nil {
				return err
			}
		}
	}

	if !modeExists {
		if err := d.electric.AddMode(mode); err != nil {
			return err
		}
	}
	return nil
}

// ensureOnlyModes adjusts the electric modes of the device such that only the given modes exists.
func (d *Device) ensureOnlyModes(modes []*traits.ElectricMode) error {
	modesToAdd := make(map[string]*traits.ElectricMode)
	modesToCreate := make([]*traits.ElectricMode, 0)
	modesToUpdate := make(map[string]*traits.ElectricMode)
	modesToDelete := make([]string, 0)

	oldModes := d.electric.Modes(nil)
	oldModesById := make(map[string]*traits.ElectricMode)
	for _, mode := range oldModes {
		oldModesById[mode.Id] = mode
	}

	for _, mode := range modes {
		if mode.Id == "" {
			modesToCreate = append(modesToCreate, mode)
		} else if _, ok := oldModesById[mode.Id]; ok {
			modesToUpdate[mode.Id] = mode
		} else {
			modesToAdd[mode.Id] = mode
		}
	}
	for _, mode := range oldModes {
		if _, ok := modesToUpdate[mode.Id]; ok {
			continue
		}
		modesToDelete = append(modesToDelete, mode.Id)
	}

	// Delete first, in case we're about to add/update a mode to contradict with an old mode.
	// Then update, for the same reason, in case a mode is update to not conflict with a future addition.
	// Add and Create can happen in any order.
	for _, id := range modesToDelete {
		if err := d.electric.DeleteMode(id); err != nil {
			return err
		}
	}
	for _, mode := range modesToUpdate {
		if _, err := d.electric.UpdateMode(mode, nil); err != nil {
			return err
		}
	}
	for _, mode := range modesToAdd {
		if err := d.electric.AddMode(mode); err != nil {
			return err
		}
	}
	for _, mode := range modesToCreate {
		if _, err := d.electric.CreateMode(mode); err != nil {
			return err
		}
	}
	return nil
}

// updateDemand sets the current electric demand of the device.
func (d *Device) updateDemand(demand float32) error {
	_, err := d.electric.UpdateDemand(&traits.ElectricDemand{
		Current: demand,
		Rating:  d.rating,
		Voltage: &d.voltage,
	})
	return err
}

// demandAt returns the electric demand at time now based on a traits.ElectricMode that activated at activated time.
// The active time is overridden by mode.StartTime if non-nil.
func (d *Device) demandAt(now, activated time.Time, mode *traits.ElectricMode) (demand float32, ok bool) {
	// todo: support ramping between segments

	cur := activated
	if mode.StartTime != nil {
		cur = mode.StartTime.AsTime()
	}
	if now.Before(cur) {
		// Can't calculate demand at times before the mode becomes active
		panic(fmt.Errorf("can't calculate demand before the mode is active: now %v < active %v", now, cur))
	}

	for _, segment := range mode.Segments {
		if segment.Length == nil {
			// infinite segment
			return d.segmentDemand(now, cur, segment), true
		}
		next := cur.Add(segment.Length.AsDuration())
		if next.After(now) {
			// active segment
			return d.segmentDemand(now, cur, segment), true
		}
		cur = next
	}
	return 0, false // no active segment found
}

// segmentDemand returns the electric demand at time now based on a single Segment that activated at active time.
func (d *Device) segmentDemand(now, active time.Time, segment *traits.ElectricMode_Segment) float32 {
	// todo: support more than just fixed segments
	_ = active
	return d.magnitudeDemand(now, segment.Magnitude)
}

func (d *Device) magnitudeDemand(now time.Time, magnitude float32) float32 {
	base := magnitude * 0.8
	return d.addNoise(base, now)
}

func (d *Device) addNoise(base float32, _ time.Time) float32 {
	// todo: add random variations to the idle current
	return base
}

func quantityDescription(p float32) traits.EnergyLevel_Quantity_Threshold {
	switch {
	case p < 0:
		return traits.EnergyLevel_Quantity_CRITICALLY_LOW
	case p == 0:
		return traits.EnergyLevel_Quantity_EMPTY
	case p <= .2:
		return traits.EnergyLevel_Quantity_LOW
	case p <= .75:
		return traits.EnergyLevel_Quantity_MEDIUM
	case p < 1:
		return traits.EnergyLevel_Quantity_HIGH
	case p == 1:
		return traits.EnergyLevel_Quantity_FULL
	case p > 1:
		return traits.EnergyLevel_Quantity_CRITICALLY_HIGH
	}
	return traits.EnergyLevel_Quantity_THRESHOLD_UNKNOWN
}
