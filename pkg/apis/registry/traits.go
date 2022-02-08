package registry

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/smart-core-os/sc-golang/pkg/trait/booking"
	"github.com/smart-core-os/sc-golang/pkg/trait/count"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/smart-core-os/sc-golang/pkg/trait/emergency"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/smart-core-os/sc-golang/pkg/trait/powersupply"
	"github.com/smart-core-os/sc-golang/pkg/trait/speaker"
	parent2 "github.com/smart-core-os/sc-playground/pkg/apis/parent"
)

// todo: Generate this file

type AirTemperatureApiRegistry struct {
	*airtemperature.ApiRouter
	parent2.Traiter
}

func (r AirTemperatureApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.AirTemperatureApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, airtemperature.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.AirTemperature)
		}
	}
	return nil
}

type AirTemperatureInfoRegistry struct{ *airtemperature.InfoRouter }

func (r AirTemperatureInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.AirTemperatureInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, airtemperature.WrapInfo(t))
		}
	}
	return nil
}

type BookingApiRegistry struct {
	*booking.ApiRouter
	parent2.Traiter
}

func (r BookingApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.BookingApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, booking.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.Booking)
		}
	}
	return nil
}

type BookingInfoRegistry struct{ *booking.InfoRouter }

func (r BookingInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.BookingInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, booking.WrapInfo(t))
		}
	}
	return nil
}

type CountApiRegistry struct {
	*count.ApiRouter
	parent2.Traiter
}

func (r CountApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.CountApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, count.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.Count)
		}
	}
	return nil
}

type CountInfoRegistry struct{ *count.InfoRouter }

func (r CountInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.CountInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, count.WrapInfo(t))
		}
	}
	return nil
}

type ElectricApiRegistry struct {
	*electric.ApiRouter
	parent2.Traiter
}

func (r ElectricApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.ElectricApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, electric.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.Electric)
		}
	}
	return nil
}

type ElectricInfoRegistry struct{ *electric.InfoRouter }

func (r ElectricInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.ElectricInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, electric.WrapInfo(t))
		}
	}
	return nil
}

type EmergencyApiRegistry struct {
	*emergency.ApiRouter
	parent2.Traiter
}

func (r EmergencyApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.EmergencyApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, emergency.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.Emergency)
		}
	}
	return nil
}

type EmergencyInfoRegistry struct{ *emergency.InfoRouter }

func (r EmergencyInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.EmergencyInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, emergency.WrapInfo(t))
		}
	}
	return nil
}

type EnergyStorageApiRegistry struct {
	*energystorage.ApiRouter
	parent2.Traiter
}

func (r EnergyStorageApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.EnergyStorageApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, energystorage.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.EnergyStorage)
		}
	}
	return nil
}

type EnergyStorageInfoRegistry struct{ *energystorage.InfoRouter }

func (r EnergyStorageInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.EnergyStorageInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, energystorage.WrapInfo(t))
		}
	}
	return nil
}

type LightApiRegistry struct {
	*light.ApiRouter
	parent2.Traiter
}

func (r LightApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.LightApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, light.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.Light)
		}
	}
	return nil
}

type LightInfoRegistry struct{ *light.InfoRouter }

func (r LightInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.LightInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, light.WrapInfo(t))
		}
	}
	return nil
}

type OccupancySensorApiRegistry struct {
	*occupancysensor.ApiRouter
	parent2.Traiter
}

func (r OccupancySensorApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.OccupancySensorApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, occupancysensor.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.OccupancySensor)
		}
	}
	return nil
}

type OccupancySensorInfoRegistry struct{ *occupancysensor.InfoRouter }

func (r OccupancySensorInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.OccupancySensorInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, occupancysensor.WrapInfo(t))
		}
	}
	return nil
}

type OnOffApiRegistry struct {
	*onoff.ApiRouter
	parent2.Traiter
}

func (r OnOffApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.OnOffApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, onoff.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.OnOff)
		}
	}
	return nil
}

type OnOffInfoRegistry struct{ *onoff.InfoRouter }

func (r OnOffInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.OnOffInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, onoff.WrapInfo(t))
		}
	}
	return nil
}

type ParentApiRegistry struct {
	*parent.ApiRouter
	parent2.Traiter
}

func (r ParentApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.ParentApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, parent.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.Parent)
		}
	}
	return nil
}

type ParentInfoRegistry struct{ *parent.InfoRouter }

func (r ParentInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.ParentInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, parent.WrapInfo(t))
		}
	}
	return nil
}

type PowerSupplyApiRegistry struct {
	*powersupply.ApiRouter
	parent2.Traiter
}

func (r PowerSupplyApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.PowerSupplyApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, powersupply.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.PowerSupply)
		}
	}
	return nil
}

type PowerSupplyInfoRegistry struct{ *powersupply.InfoRouter }

func (r PowerSupplyInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.PowerSupplyInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, powersupply.WrapInfo(t))
		}
	}
	return nil
}

type SpeakerApiRegistry struct {
	*speaker.ApiRouter
	parent2.Traiter
}

func (r SpeakerApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.SpeakerApiServer); ok {
		if r.ApiRouter != nil {
			r.Add(name, speaker.WrapApi(t))
		}
		if r.Traiter != nil {
			r.Trait(name, trait.Speaker)
		}
	}
	return nil
}

type SpeakerInfoRegistry struct{ *speaker.InfoRouter }

func (r SpeakerInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.SpeakerInfoServer); ok {
		if r.InfoRouter != nil {
			r.Add(name, speaker.WrapInfo(t))
		}
	}
	return nil
}
