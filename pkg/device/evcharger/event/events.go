package event

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-playground/pkg/device"
)

// PlugIn is a Named Event that represents a device plugging in to a charger.
type PlugIn struct {
	device.Event
	device.Named
	// Level encodes the current charge level when the vehicle was plugged in.
	// todo: the virtual device doesn't take this into account yet
	Level *traits.EnergyLevel_Quantity
	// Full describes the values the plugged in device would have when full.
	//
	// Example
	//
	//    Percentage: 100, // when full, the ev is at 100%
	//    DistanceKm: 300, // when full, the ev can travel 300 km
	//    EnergyKwh:  100, // when full, the ev hold 100 kWh of energy
	Full *traits.EnergyLevel_Quantity
	// Supported charging modes of the device that was just plugged in.
	// Empty implies the charger can make the decision.
	// The first mode is the preferred mode and will be marked as the Normal mode.
	Modes []*ChargeMode
}

// Unplug is a Named Event that represents a device unplugging from a charger.
type Unplug struct {
	device.Event
	device.Named
}

// ChargeStart is a Named Event that represents the moment charging starts.
type ChargeStart struct {
	device.Event
	device.Named
}

// ModeChange is a Named Event that represents when the charging mode has been asked to change.
type ModeChange struct {
	device.Event
	device.Named
	Id string
}

// ChargeMode describes a mode of charging, similar to traits.ElectricMode but more focused.
type ChargeMode struct {
	Id, Title, Description string

	Segments []*traits.ElectricMode_Segment
}
