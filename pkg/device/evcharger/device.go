package evcharger

import (
	"time"

	"github.com/smart-core-os/sc-playground/pkg/apis/registry"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
	"github.com/smart-core-os/sc-playground/internal/util/errs"
	"github.com/smart-core-os/sc-playground/pkg/device/evcharger/event"
	"github.com/smart-core-os/sc-playground/pkg/timeline"
	"github.com/smart-core-os/sc-playground/pkg/timeline/tlutil"
)

// Device is a virtual electric vehicle charger.
type Device struct {
	name string

	rating  float32
	voltage float32

	model         *event.Model
	electric      *electric.Model
	energyStorage *energystorage.Model

	electricApi      traits.ElectricApiServer
	energyStorageApi traits.EnergyStorageApiServer
}

// New creates a new Device.
func New(name string, tl timeline.TL, opts ...Opt) *Device {
	setup := calcSetup(opts...)

	if !setup.dedicatedTL {
		tl = tlutil.FilterByName(tl, name)
	}

	electricModel := electric.NewModel(setup.clock)
	energyStorageModel := energystorage.NewModel()

	d := &Device{
		name:          name,
		model:         &event.Model{TL: tl},
		electric:      electricModel,
		energyStorage: energyStorageModel,

		electricApi:      electric.NewModelServer(electricModel),
		energyStorageApi: energystorage.NewModelServer(energyStorageModel),

		rating:  setup.rating,
		voltage: setup.voltage,
	}

	if setup.dispatcher != nil {
		d.electricApi = event.CaptureElectricInput(d.electricApi, setup.clock, setup.dispatcher)
		d.energyStorageApi = event.CaptureEnergyStorageInput(d.energyStorageApi, setup.clock, setup.dispatcher)
	}

	return d
}

func (d *Device) Scrub(t time.Time) error {
	m := d.model
	if err := m.Scrub(t); err != nil {
		return err
	}

	switch {
	case m.IsIdle():
		return d.setIdle()
	case m.IsPluggedIn() && !m.IsCharging():
		return d.setPluggedIn()
	default:
		return d.setCharging()
	}
}

func (d *Device) Publish(reg registry.Registry) (err error) {
	return errs.First(
		reg.Register(d.name, d.electricApi),
		reg.Register(d.name, d.energyStorageApi),
	)
}
