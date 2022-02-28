package evcharger

import (
	"time"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadata"
	"github.com/smart-core-os/sc-playground/pkg/node"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
	"github.com/smart-core-os/sc-playground/pkg/device/evcharger/event"
	"github.com/smart-core-os/sc-playground/pkg/timeline"
	"github.com/smart-core-os/sc-playground/pkg/timeline/tlutil"
)

func Activate(n *node.Node) {
	n.AddRouter(NewApiRouter())
}

// Device is a virtual electric vehicle charger.
type Device struct {
	name string

	rating  float32
	voltage float32

	model         *event.Model
	electric      *electric.Model
	energyStorage *energystorage.Model
	metadata      *metadata.Model

	electricApi      traits.ElectricApiServer
	energyStorageApi traits.EnergyStorageApiServer
	metadataApi      traits.MetadataApiServer

	api *ApiServer
}

// New creates a new Device.
func New(name string, tl timeline.TL, opts ...Opt) *Device {
	setup := calcSetup(opts...)

	if !setup.dedicatedTL {
		tl = tlutil.FilterByName(tl, name)
	}

	electricModel := electric.NewModel(setup.clock)
	energyStorageModel := energystorage.NewModel()
	metadataModel := metadata.NewModel(resource.WithInitialValue(newMetadata(name)))

	d := &Device{
		name:          name,
		model:         &event.Model{TL: tl},
		electric:      electricModel,
		energyStorage: energyStorageModel,
		metadata:      metadataModel,

		electricApi:      electric.NewModelServer(electricModel),
		energyStorageApi: energystorage.NewModelServer(energyStorageModel),
		metadataApi:      metadata.NewModelServer(metadataModel),

		api: &ApiServer{clock: clock.Real(), dispatcher: setup.dispatcher},

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

func (d *Device) Publish(announcer node.Announcer) {
	announcer.Announce(
		d.name,
		node.HasClient(WrapApi(d.api)),
		node.HasSimulation(d),
		node.HasTrait(trait.Electric, node.WithClients(electric.WrapApi(d.electricApi)), node.NoAddMetadata()),
		node.HasTrait(trait.EnergyStorage, node.WithClients(energystorage.WrapApi(d.energyStorageApi)), node.NoAddMetadata()),
		node.HasTrait(trait.Metadata, node.WithClients(metadata.WrapApi(d.metadataApi)), node.NoAddMetadata()),
	)
}
