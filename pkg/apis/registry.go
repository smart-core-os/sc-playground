package apis

import (
	"github.com/smart-core-os/sc-api/go/traits"
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
)

type Publisher interface {
	Publish(reg Registry) error
}

type Registry interface {
	Register(name string, impl interface{}) error
}

type RegistrySlice []Registry

func (r RegistrySlice) Register(name string, impl interface{}) (err error) {
	for _, reg := range r {
		err = reg.Register(name, impl)
	}
	return err
}

type AirTemperatureApiRegistry struct{ *airtemperature.ApiRouter }

func (r AirTemperatureApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.AirTemperatureApiServer); ok {
		r.Add(name, airtemperature.WrapApi(t))
	}
	return nil
}

type AirTemperatureInfoRegistry struct{ *airtemperature.InfoRouter }

func (r AirTemperatureInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.AirTemperatureInfoServer); ok {
		r.Add(name, airtemperature.WrapInfo(t))
	}
	return nil
}

type BookingApiRegistry struct{ *booking.ApiRouter }

func (r BookingApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.BookingApiServer); ok {
		r.Add(name, booking.WrapApi(t))
	}
	return nil
}

type BookingInfoRegistry struct{ *booking.InfoRouter }

func (r BookingInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.BookingInfoServer); ok {
		r.Add(name, booking.WrapInfo(t))
	}
	return nil
}

type CountApiRegistry struct{ *count.ApiRouter }

func (r CountApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.CountApiServer); ok {
		r.Add(name, count.WrapApi(t))
	}
	return nil
}

type CountInfoRegistry struct{ *count.InfoRouter }

func (r CountInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.CountInfoServer); ok {
		r.Add(name, count.WrapInfo(t))
	}
	return nil
}

type ElectricApiRegistry struct{ *electric.ApiRouter }

func (r ElectricApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.ElectricApiServer); ok {
		r.Add(name, electric.WrapApi(t))
	}
	return nil
}

type ElectricInfoRegistry struct{ *electric.InfoRouter }

func (r ElectricInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.ElectricInfoServer); ok {
		r.Add(name, electric.WrapInfo(t))
	}
	return nil
}

type EmergencyApiRegistry struct{ *emergency.ApiRouter }

func (r EmergencyApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.EmergencyApiServer); ok {
		r.Add(name, emergency.WrapApi(t))
	}
	return nil
}

type EmergencyInfoRegistry struct{ *emergency.InfoRouter }

func (r EmergencyInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.EmergencyInfoServer); ok {
		r.Add(name, emergency.WrapInfo(t))
	}
	return nil
}

type EnergyStorageApiRegistry struct{ *energystorage.ApiRouter }

func (r EnergyStorageApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.EnergyStorageApiServer); ok {
		r.Add(name, energystorage.WrapApi(t))
	}
	return nil
}

type EnergyStorageInfoRegistry struct{ *energystorage.InfoRouter }

func (r EnergyStorageInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.EnergyStorageInfoServer); ok {
		r.Add(name, energystorage.WrapInfo(t))
	}
	return nil
}

type LightApiRegistry struct{ *light.ApiRouter }

func (r LightApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.LightApiServer); ok {
		r.Add(name, light.WrapApi(t))
	}
	return nil
}

type LightInfoRegistry struct{ *light.InfoRouter }

func (r LightInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.LightInfoServer); ok {
		r.Add(name, light.WrapInfo(t))
	}
	return nil
}

type OccupancySensorApiRegistry struct{ *occupancysensor.ApiRouter }

func (r OccupancySensorApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.OccupancySensorApiServer); ok {
		r.Add(name, occupancysensor.WrapApi(t))
	}
	return nil
}

type OccupancySensorInfoRegistry struct{ *occupancysensor.InfoRouter }

func (r OccupancySensorInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.OccupancySensorInfoServer); ok {
		r.Add(name, occupancysensor.WrapInfo(t))
	}
	return nil
}

type OnOffApiRegistry struct{ *onoff.ApiRouter }

func (r OnOffApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.OnOffApiServer); ok {
		r.Add(name, onoff.WrapApi(t))
	}
	return nil
}

type OnOffInfoRegistry struct{ *onoff.InfoRouter }

func (r OnOffInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.OnOffInfoServer); ok {
		r.Add(name, onoff.WrapInfo(t))
	}
	return nil
}

type ParentApiRegistry struct{ *parent.ApiRouter }

func (r ParentApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.ParentApiServer); ok {
		r.Add(name, parent.WrapApi(t))
	}
	return nil
}

type ParentInfoRegistry struct{ *parent.InfoRouter }

func (r ParentInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.ParentInfoServer); ok {
		r.Add(name, parent.WrapInfo(t))
	}
	return nil
}

type PowerSupplyApiRegistry struct{ *powersupply.ApiRouter }

func (r PowerSupplyApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.PowerSupplyApiServer); ok {
		r.Add(name, powersupply.WrapApi(t))
	}
	return nil
}

type PowerSupplyInfoRegistry struct{ *powersupply.InfoRouter }

func (r PowerSupplyInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.PowerSupplyInfoServer); ok {
		r.Add(name, powersupply.WrapInfo(t))
	}
	return nil
}

type SpeakerApiRegistry struct{ *speaker.ApiRouter }

func (r SpeakerApiRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.SpeakerApiServer); ok {
		r.Add(name, speaker.WrapApi(t))
	}
	return nil
}

type SpeakerInfoRegistry struct{ *speaker.InfoRouter }

func (r SpeakerInfoRegistry) Register(name string, impl interface{}) error {
	if t, ok := impl.(traits.SpeakerInfoServer); ok {
		r.Add(name, speaker.WrapInfo(t))
	}
	return nil
}
