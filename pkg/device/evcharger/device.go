package evcharger

import (
	"time"

	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
	"github.com/smart-core-os/sc-playground/internal/util/errs"
	"github.com/smart-core-os/sc-playground/pkg/apis"
	"github.com/smart-core-os/sc-playground/pkg/device/evcharger/event"
	"github.com/smart-core-os/sc-playground/pkg/timeline"
	"github.com/smart-core-os/sc-playground/pkg/timeline/tlutil"
)

// Device is a virtual electric vehicle charger.
type Device struct {
	name string

	rating  float32
	voltage float32

	electric      *electric.Model
	energyStorage *energystorage.Model

	model *event.Model
}

// New creates a new Device.
func New(name string, tl timeline.TL, opts ...Opt) *Device {
	// option processing
	setup := &opt{}
	for _, opt := range DefaultOpts {
		opt(setup)
	}
	for _, opt := range opts {
		opt(setup)
	}

	if !setup.dedicatedTL {
		tl = tlutil.FilterByName(tl, name)
	}

	electricModel := electric.NewModel(clock.Real())
	energyStorageModel := energystorage.NewModel()
	d := &Device{
		name:          name,
		electric:      electricModel,
		energyStorage: energyStorageModel,
		model:         &event.Model{TL: tl},

		rating:  setup.rating,
		voltage: setup.voltage,
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

func (d *Device) Publish(reg apis.Registry) (err error) {
	return errs.First(
		reg.Register(d.name, electric.NewModelServer(d.electric)),
		reg.Register(d.name, energystorage.NewModelServer(d.energyStorage)),
	)
}
